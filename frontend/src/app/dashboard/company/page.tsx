'use client';
import React, { useState } from 'react';
import Link from 'next/link';
import { PrimaryButton } from '@/components/buttons/PrimaryButton';
import { useAuthCheck } from '@/hooks/useAuthCheck';
import {
  companyScrape,
  generateReasons,
  saveCompanyReason,
} from '@/lib/api/company';
import { MatchItem } from '@/lib/types';
import { errorToast, successToast } from '@/lib/toast';

const defaultQuestions = {
  reasonInterest: '',
  attractiveService: '',
  relatedExperience: '',
};

export default function Company() {
  useAuthCheck();
  const [companyUrl, setCompanyUrl] = useState('');
  const [matchResult, setMatchResult] = useState<MatchItem[]>([]);
  const [isScraping, setIsScraping] = useState(false);
  const [isGenerating, setIsGenerating] = useState(false);
  const [showQuestions, setShowQuestions] = useState(false);
  const [questions, setQuestions] = useState(defaultQuestions);
  const [companyReasons, setCompanyReasons] = useState('');
  const [companyName, setCompanyName] = useState('');
  const onChangeCompanyUrl = (event: React.ChangeEvent<HTMLInputElement>) =>
    setCompanyUrl(event.target.value);

  const onSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    submit();
  };

  const submit = async () => {
    if (!companyUrl.trim()) {
      errorToast('企業HPを入力してください。');
      return;
    }

    if (!companyName.trim()) {
      errorToast('企業名を入力してください。');
      return;
    }

    let url;
    try {
      url = new URL(companyUrl);
    } catch {
      errorToast('有効なURLを入力してください。');
      return;
    }

    if (url.protocol !== 'http:' && url.protocol !== 'https:') {
      errorToast('URLは http または https で始まる必要があります。');
      return;
    }
    setIsScraping(true);
    try {
      const res = await companyScrape({ companyUrl });
      const result = res.matchResult || [];
      if (result.length === 0) {
        errorToast('企業情報の取得に失敗しました。');
        return;
      }
      setMatchResult(res.matchResult || '');
    } catch (e) {
      console.error(e);
    } finally {
      setIsScraping(false);
    }
  };

  const onClickGenerateReasons = async () => {
    setIsGenerating(true);
    try {
      const res = await generateReasons({ matchResult, questions });
      setCompanyReasons(res.reason);
    } catch (e) {
      console.error('志望理由の生成に失敗:', e);
    } finally {
      setIsGenerating(false);
    }
  };

  const onClickPostReasons = async () => {
    try {
      await saveCompanyReason({
        content: companyReasons,
        companyName,
        companyUrl,
      });
      successToast('保存に成功しました。');
    } catch (e) {
      console.error('志望理由の保存に失敗:', e);
      errorToast('保存に失敗しました。');
    }
  };

  return (
    <div className='max-w-3xl mx-auto px-4 py-8'>
      <Link
        href='/dashboard'
        className='text-blue-600 underline hover:text-blue-400'
      >
        ← ダッシュボードへ戻る
      </Link>

      <h1 className='text-2xl font-bold text-blue-700 mt-6 mb-4'>
        企業分析ツール
      </h1>

      <form
        onSubmit={onSubmit}
        className='flex flex-col sm:flex-row items-start gap-4 mb-6'
      >
        <input
          type='text'
          value={companyUrl}
          onChange={onChangeCompanyUrl}
          placeholder='企業のURLを入力してください'
          className='w-full sm:flex-1 px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-400'
        />
        <input
          type='text'
          value={companyName}
          onChange={(e) => setCompanyName(e.target.value)}
          placeholder='企業名を入力してください'
          className='w-full sm:flex-1 px-4 py-2 border border-gray-300 rounded-md shadow-sm'
        />
        <PrimaryButton type='submit' disabled={isScraping}>
          {isScraping ? (
            <span className='animate-pulse'>取得中...</span>
          ) : (
            '企業情報取得'
          )}
        </PrimaryButton>
      </form>
      <div>
        <p className='font-semibold text-gray-600 mb-2'>分析結果:</p>
        <div className='space-y-4'>
          {matchResult.map((item, index) => (
            <div
              key={index}
              className='p-4 border rounded-md shadow-sm bg-white'
            >
              <h3 className='font-bold text-blue-700 underline'>
                {item.axis}:{'★'.repeat(item.score)}
                {'☆'.repeat(5 - item.score)}
              </h3>
              <p className='text-gray-700 mt-1'>{item.reason}</p>
            </div>
          ))}
        </div>
      </div>
      {matchResult.length > 0 && (
        <div className='mt-5 flex justify-center'>
          <PrimaryButton onClick={() => setShowQuestions(true)}>
            志望理由を作成
          </PrimaryButton>
        </div>
      )}

      {showQuestions && (
        <form className='space-y-6 mt-8'>
          <div>
            <label>この企業に興味を持った理由は何ですか？</label>
            <textarea
              name='reasonInterest'
              value={questions.reasonInterest}
              onChange={(e) =>
                setQuestions({ ...questions, reasonInterest: e.target.value })
              }
              className='w-full border p-2 rounded-md'
            />
          </div>
          <div>
            <label>この企業のどんな製品・サービスに魅力を感じましたか？</label>
            <textarea
              name='attractiveService'
              value={questions.attractiveService}
              onChange={(e) =>
                setQuestions({
                  ...questions,
                  attractiveService: e.target.value,
                })
              }
              className='w-full border p-2 rounded-md'
            />
          </div>
          <div>
            <label>
              その製品・サービスに関して、あなたの過去経験と結びつくものはありますか？
            </label>
            <textarea
              name='relatedExperience'
              value={questions.relatedExperience}
              onChange={(e) =>
                setQuestions({
                  ...questions,
                  relatedExperience: e.target.value,
                })
              }
              className='w-full border p-2 rounded-md'
            />
          </div>
          <div className='text-right'>
            <PrimaryButton
              onClick={onClickGenerateReasons}
              disabled={isGenerating}
            >
              {isGenerating ? (
                <span className='animate-pulse'>取得中...</span>
              ) : (
                '志望理由を生成'
              )}
            </PrimaryButton>
          </div>
        </form>
      )}
      {companyReasons && (
        <div className='mt-8 p-4 border rounded-md bg-gray-50'>
          <h2 className='text-lg font-semibold text-blue-700 mb-2'>
            生成された志望理由（編集可）:
          </h2>
          <textarea
            className='w-full border rounded-md p-2 text-gray-700'
            value={companyReasons}
            onChange={(e) => setCompanyReasons(e.target.value)}
            rows={6}
          />
          <div className='text-right mt-2'>
            <PrimaryButton onClick={onClickPostReasons}>保存</PrimaryButton>
          </div>
        </div>
      )}
    </div>
  );
}
