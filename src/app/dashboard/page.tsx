'use client';

import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { jwtDecode } from 'jwt-decode';
import { User, Building, ArrowLeft } from 'lucide-react';
import { SkeltonButton } from '@/components/buttons/Skeltonbutton';
import Link from 'next/link';

export default function Dashboard() {
  const router = useRouter();
  const [email, setEmail] = useState<string>('');

  useEffect(() => {
    const token = localStorage.getItem('access_token');
    if (token && token.split('.').length === 3) {
      const decodedUser = jwtDecode<{
        id: number;
        email: string;
        iat: number;
        exp: number;
      }>(token);
      setEmail(decodedUser.email);
    } else if (!token) {
      router.push('/');
    }
  }, [router]);

  return (
    <div className='py-8 px-20'>
      <Link
        href='/'
        className='flex items-center gap-2 text-blue-600 mb-6'
      >
        <ArrowLeft />
        <span>TOPに戻る</span>
      </Link>
      <h1 className='text-3xl font-bold mb-8'>ダッシュボード</h1>
      <div className='bg-white rounded-lg shadow-md p-6 mb-8'>
        <div className='flex items-center gap-4 mb-4'>
          <div className='bg-blue-100 p-3 rounded-full'>
            <User className='h-6 w-6 text-blue-600' />
          </div>
          <div>
            <h2 className='text-xl font-semibold'>プロフィール</h2>
            <p className='text-gray-600'>{email}</p>
          </div>
        </div>

        <SkeltonButton href='/profile/edit'>プロフィールを編集</SkeltonButton>
      </div>
      <div className='bg-white rounded-lg shadow-md p-6 mb-8'>
        <div className='flex items-center justify-between'>
          <div className='flex items-center gap-4'>
            <div className='bg-blue-100 p-3 rounded-full'>
              <Building className='h-6 w-6 text-blue-600' />
            </div>
            <h2 className='text-xl font-semibold'>企業分析</h2>
          </div>
          <SkeltonButton href='/dashboard/company'>詳細へ</SkeltonButton>
        </div>
      </div>
      <h2 className='text-2xl font-bold mb-4'>志望企業</h2>
    </div>
  );
}
