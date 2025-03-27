import { NextRequest, NextResponse } from 'next/server';
import puppeteer from 'puppeteer';
import { load } from 'cheerio';

export async function POST(req: NextRequest) {
  try {
    const { companyUrl } = await req.json();

    if (!companyUrl) {
      return NextResponse.json(
        { success: false, error: 'No companyUrl provided' },
        { status: 400 },
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
        try {
          const fullUrl = new URL(href, companyUrl).toString();
          linkSet.add(fullUrl);
        } catch (_) {
          // console.warn(`Invalid URL skipped: ${href}`);
        }
      }
    });

    const subPageResults: { url: string; text: string }[] = [];

    for (const url of linkSet) {
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

    return NextResponse.json({
      success: true,
      subPages: subPageResults,
    });
  } catch (e) {
    return NextResponse.json(e);
  }
}
