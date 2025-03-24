'use client';
import React, { useState } from 'react';
import Link from 'next/link';
import { PrimaryButton } from '@/components/buttons/PrimaryButton';
import { signUp } from '@/lib/api/auth';
import { errorToast, successToast } from '@/lib/toast';
import { useRouter } from 'next/navigation';

export default function SignIn() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const onChangeEmail = (event: React.ChangeEvent<HTMLInputElement>) =>
    setEmail(event.target.value);
  const onChangePassword = (event: React.ChangeEvent<HTMLInputElement>) =>
    setPassword(event.target.value);

  const onSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    submit();
  };
  const submit = async () => {
    const res = await signUp({ email, password });
    const token = res?.token;

    if (token) {
      localStorage.setItem('access_token', token);
      successToast('新規登録に成功しました。');
      router.push('/user/dashboard');
    } else {
      errorToast('そのメールアドレスは既に使用されています。');
    }
  };

  return (
    <div className='flex flex-col items-center justify-center min-h-screen bg-gray-100 p-6'>
      <div className='bg-white shadow-lg rounded-2xl p-8 max-w-md w-full text-center'>
        <h1 className='text-3xl font-bold text-gray-800 mb-4'>
          サインアップ画面
        </h1>
        <form onSubmit={onSubmit} className='mt-6 flex flex-col gap-4'>
          <input
            onChange={onChangeEmail}
            type='email'
            placeholder='メールアドレス'
            className='w-full px-4 py-3 border border-gray-300 rounded-lg text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500'
          />
          <input
            onChange={onChangePassword}
            type='password'
            placeholder='パスワード'
            className='w-full px-4 py-3 border border-gray-300 rounded-lg text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500'
          />
          <PrimaryButton disabled={!email || !password}>新規登録</PrimaryButton>
        </form>
        <Link href='/'>TOPへ</Link>
      </div>
    </div>
  );
}
