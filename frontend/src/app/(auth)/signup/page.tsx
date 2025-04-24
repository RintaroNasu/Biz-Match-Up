'use client';
import React, { useState } from 'react';
import Link from 'next/link';
import { PrimaryButton } from '../../../components/buttons/PrimaryButton';
import { signUp } from '../../../lib/api/auth';
import { errorToast, successToast } from '../../../lib/toast';
import { useRouter } from 'next/navigation';

export default function SignUp() {
  const router = useRouter();
  const [form, setForm] = useState({
    email: '',
    password: '',
    name: '',
    desiredJobType: '',
    desiredLocation: '',
    desiredCompanySize: '',
    careerAxis1: '',
    careerAxis2: '',
    selfPr: '',
  });

  const onChangeForm = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    console.log(form);
    const res = await signUp(form);
    const token = res?.token;

    if (token) {
      localStorage.setItem('access_token', token);
      successToast('新規登録に成功しました。');
      router.push('/dashboard');
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
            name='email'
            type='email'
            placeholder='メールアドレス'
            value={form.email}
            onChange={onChangeForm}
            className='input'
          />
          <input
            name='password'
            type='password'
            placeholder='パスワード'
            value={form.password}
            onChange={onChangeForm}
            className='input'
          />
          <input
            name='name'
            placeholder='名前'
            value={form.name}
            onChange={onChangeForm}
            className='input'
          />
          <input
            name='desiredJobType'
            placeholder='志望職種'
            value={form.desiredJobType}
            onChange={onChangeForm}
            className='input'
          />
          <input
            name='desiredLocation'
            placeholder='志望勤務地'
            value={form.desiredLocation}
            onChange={onChangeForm}
            className='input'
          />
          <input
            name='desiredCompanySize'
            placeholder='志望企業の規模'
            value={form.desiredCompanySize}
            onChange={onChangeForm}
            className='input'
          />
          <input
            name='careerAxis1'
            placeholder='就活軸①（任意）'
            value={form.careerAxis1}
            onChange={onChangeForm}
            className='input'
          />
          <input
            name='careerAxis2'
            placeholder='就活軸②（任意）'
            value={form.careerAxis2}
            onChange={onChangeForm}
            className='input'
          />
          <textarea
            name='selfPr'
            placeholder='自己PR（任意）'
            value={form.selfPr}
            onChange={onChangeForm}
            className='input h-24'
          />
          <PrimaryButton disabled={!form.email || !form.password || !form.name}>
            新規登録
          </PrimaryButton>
        </form>
        <Link href='/' className='mt-4 block text-blue-500'>
          TOPへ
        </Link>
      </div>
    </div>
  );
}
