'use client';
import React, { useCallback } from 'react';
import Link from 'next/link';
import { SkeltonButton } from '../components/buttons/SkeltonButton';
import { useRouter } from 'next/navigation';
import { successToast } from '../lib/toast';
import { PrimaryButton } from '../components/buttons/PrimaryButton';
import { useAuthCheck } from '@/hooks/useAuthCheck';

export default function Home() {
  const { user } = useAuthCheck();
  const router = useRouter();
  const isLoggedIn = !!user;

  const onClickLogout = useCallback(() => {
    localStorage.removeItem('access_token');
    successToast('ログアウトしました。');
    router.push('/');
  }, [router]);

  return (
    <div className='flex flex-col items-center justify-center min-h-screen bg-gray-100 p-6'>
      <div className='bg-white shadow-lg rounded-2xl p-8 max-w-md text-center'>
        <h1 className='text-3xl font-bold text-gray-800 mb-4'>BizMatchUp</h1>
        <p className='text-lg text-gray-600'>アプリ概要について</p>
        <p className='mt-2 text-gray-500'>
          企業研究から志望理由作成までを自動化させるアプリになります。
        </p>
        <div className='mt-6 flex flex-col gap-4'>
          {!isLoggedIn && (
            <>
              <Link
                href='/signup'
                className='w-full px-6 py-3 text-white bg-blue-600 hover:bg-blue-700 rounded-lg text-lg font-semibold shadow-md transition'
              >
                新規登録
              </Link>
              <SkeltonButton href='/signin'>ログイン</SkeltonButton>
            </>
          )}

          {isLoggedIn && (
            <>
              <SkeltonButton href='/dashboard'>ダッシュボード</SkeltonButton>
              <PrimaryButton onClick={onClickLogout}>ログアウト</PrimaryButton>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
