import { NextRequest, NextResponse } from 'next/server';
import puppeteer from 'puppeteer';
import { load } from 'cheerio';
import { OpenAI } from 'openai';
import jwt from 'jsonwebtoken';
import prisma from '@/lib/prisma';

const openai = new OpenAI({ apiKey: process.env.OPENAI_API_KEY! });

export async function POST(req: NextRequest) {
  try {
    const { companyUrl } = await req.json();
    const token = req.headers.get('authorization')?.split('Bearer ')[1];
    console.log('token:', token);
    if (!companyUrl || !token) {
      return NextResponse.json(
        { success: false, error: 'Missing companyUrl or token' },
        { status: 400 },
      );
    }

    const decoded = jwt.verify(token, process.env.JWT_SECRET!) as {
      id: number;
    };
    console.log('decoded:', decoded);
    const user = await prisma.user.findUnique({ where: { id: decoded.id } });
    console.log('user:', user);
    if (!user) {
      return NextResponse.json(
        { success: false, error: 'User not found' },
        { status: 404 },
      );
    }

    const browser = await puppeteer.launch({
      headless: true,
      args: ['--no-sandbox', '--disable-setuid-sandbox'],
    });

    const page = await browser.newPage();
    await page.goto(companyUrl, { waitUntil: 'networkidle0', timeout: 90000 });

    const html = await page.content();
    const $ = load(html);

    const linkSet = new Set<string>();

    $('header a').each((_, el) => {
      const href = $(el).attr('href');
      if (href && !href.startsWith('javascript:')) {
        const fullUrl = new URL(href, companyUrl).toString();
        linkSet.add(fullUrl);
      }
    });

    const importantPages = Array.from(linkSet)
      .filter((url) =>
        /recruit|job|vision|engineer|business|services|technology|interview|message|culture|about|company|profile|corporate|access|location|history/.test(
          url,
        ),
      )
      .slice(0, 5);

    const subPageResults: { url: string; text: string }[] = [];

    for (const url of importantPages) {
      try {
        const subPage = await browser.newPage();
        await subPage.goto(url, { waitUntil: 'networkidle2', timeout: 60000 });

        const subHtml = await subPage.content();
        const $sub = load(subHtml);
        const text = $sub('body')
          .text()
          .replace(/\s+/g, ' ')
          .trim()
          .slice(0, 3000);

        subPageResults.push({ url, text });
        await subPage.close();
      } catch (e) {
        return NextResponse.json(e);
      }
    }
    await browser.close();
    console.log('subPageResults:', subPageResults);
    const formattedPages = subPageResults
      .slice(0, 5)
      .map(({ url, text }) => {
        const cleanText = text
          .replace(/<iframe[\s\S]*?<\/iframe>/gi, '')
          .replace(/<[^>]*>/g, '')
          .replace(/\s+/g, ' ')
          .trim()
          .slice(0, 1000);

        return `▼ URL: ${url}\n${cleanText}`;
      })
      .join('\n\n');
    const prompt = `
      以下のユーザープロフィールと企業のWebページ情報をもとに、以下の評価軸でそれぞれマッチ度（★1〜5）を出してください。

      【評価軸】
      1. 専門領域での開発
      2. 裁量権・自由度
      3. 勤務地
      4. 会社の規模

      【ユーザー情報】
      名前: ${user.name}
      志望職種: ${user.desiredJobType}
      志望勤務地: ${user.desiredLocation}
      志望企業の規模: ${user.desiredCompanySize}
      就活軸①: ${user.careerAxis1}
      就活軸②: ${user.careerAxis2}
      自己PR: ${user.selfPr}

      【企業情報（Webページから抽出）】
      ${formattedPages}
      各評価軸ごとに、以下の形式でマッチ度を出力してください：

      例：
      1. 専門領域での開発: ★★★★☆
      理由: 〇〇
      2. 裁量権・自由度: ★★☆☆☆
      理由: △△ 
      ...`;

    const chatRes = await openai.chat.completions.create({
      messages: [{ role: 'user', content: prompt }],
      model: 'gpt-4',
    });

    const result = chatRes.choices[0].message.content;
    console.log('result:', result);
    return NextResponse.json({ success: true, matchResult: result });
  } catch (error: any) {
    console.error('Error:', error);
    return NextResponse.json(
      { success: false, error: error.message },
      { status: 500 },
    );
  }
}
