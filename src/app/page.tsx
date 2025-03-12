import React from 'react';
import Link from 'next/link';

export default function Home() {
  return (
    <div className='flex flex-col items-center justify-center min-h-screen bg-gray-100 p-6'>
      <div className='bg-white shadow-lg rounded-2xl p-8 max-w-md text-center'>
        <h1 className='text-3xl font-bold text-gray-800 mb-4'>BizMatchUp</h1>
        <p className='text-lg text-gray-600'>アプリ概要について</p>
        <p className='mt-2 text-gray-500'>
          企業研究から志望理由作成までを自動化させるアプリになります。
        </p>
        <div className='mt-6 flex flex-col gap-4'>
          <Link
            href='/signup'
            className='w-full px-6 py-3 text-white bg-blue-600 hover:bg-blue-700 rounded-lg text-lg font-semibold shadow-md transition'
          >
            新規登録
          </Link>
          <Link
            href='/login'
            className='w-full px-6 py-3 text-blue-600 border border-blue-600 hover:bg-blue-100 rounded-lg text-lg font-semibold shadow-md transition'
          >
            ログイン
          </Link>
        </div>
      </div>
    </div>
  );
}
