'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { User, Building, ArrowLeft } from 'lucide-react';
import { SkeltonButton } from '@/components/buttons/Skeltonbutton';
import { useAuthCheck } from '@/hooks/useAuthCheck';
import { getCompanyReasons } from '@/lib/api/company';

type Reason = {
  ID: number;
  Content: string;
  CreatedAt: string;
  CompanyName: string;
  CompanyUrl: string;
};

export default function Dashboard() {
  const { user } = useAuthCheck();
  const [reasons, setReasons] = useState<Reason[]>([]);

  useEffect(() => {
    const fetchReasons = async () => {
      try {
        const res = await getCompanyReasons();
        setReasons(res);
      } catch (error) {
        console.error('理由取得失敗:', error);
      }
    };
    if (user) {
      fetchReasons();
    }
  }, [user]);

  return (
    <div className='py-8 px-20'>
      <Link href='/' className='flex items-center gap-2 text-blue-600 mb-6'>
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
            <p className='text-gray-600'>{user?.email}</p>
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
      <div className='bg-white rounded-lg shadow-md p-6'>
        {reasons.map((reason) => (
          <div key={reason.ID} className='mb-4 p-4 border rounded-md bg-white'>
            <p className='text-gray-800 font-semibold'>{reason.CompanyName}</p>
            <p className='text-sm text-blue-600 break-all mb-2'>
              <a
                href={reason.CompanyUrl}
                target='_blank'
                rel='noopener noreferrer'
              >
                {reason.CompanyUrl}
              </a>
            </p>
            <p className='text-gray-800'>{reason.Content}</p>
            <p className='text-sm text-gray-500 mt-2'>
              作成日: {new Date(reason.CreatedAt).toLocaleDateString()}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
