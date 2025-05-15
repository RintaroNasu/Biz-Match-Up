'use client';
import React, { useState } from 'react';
import Link from 'next/link';
import { PrimaryButton } from '../../../components/buttons/PrimaryButton';
import { companyScrape } from '../../../lib/api/companyScrape';

export default function Company() {
  const [companyUrl, setCompanyUrl] = useState('');
  const [matchResult, setMatchResult] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const onChangeCompanyUrl = (event: React.ChangeEvent<HTMLInputElement>) =>
    setCompanyUrl(event.target.value);

  const onSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    submit();
  };

  const submit = async () => {
    setIsLoading(true);
    try {
      const res = await companyScrape({ companyUrl });
      console.log('res', res);
      setMatchResult(res.matchResult || '');
    } catch (e) {
      console.error(e);
    } finally {
      setIsLoading(false);
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
        <PrimaryButton disabled={isLoading}>
          {isLoading ? (
            <span className='animate-pulse'>取得中...</span>
          ) : (
            '企業情報取得'
          )}
        </PrimaryButton>
      </form>
      <div>
        <p className='font-semibold text-gray-600 mb-2'>分析結果:</p>
        {matchResult ? (
          <pre className='whitespace-pre-wrap text-sm text-gray-800'>
            {matchResult}
          </pre>
        ) : (
          <p className='text-sm text-gray-500'>まだ情報がありません。</p>
        )}
      </div>
    </div>
  );
}
