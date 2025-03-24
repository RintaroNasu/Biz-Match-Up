'use client';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { jwtDecode } from 'jwt-decode';
import { SkeltonButton } from '@/components/buttons/Skeltonbutton';

type DecodedUser = {
  email: string;
  iat: number;
  exp: number;
};

export default function Dashboard() {
  const router = useRouter();
  const [email, setEmail] = useState<string>();

  useEffect(() => {
    const token = localStorage.getItem('access_token');
    if (token && token.split('.').length === 3) {
      const decodedUser = jwtDecode<DecodedUser>(token);
      setEmail(decodedUser.email);
    } else if (!token) {
      router.push('/');
    }
  }, [router]);

  return (
    <div className='flex flex-col items-center mt-10 gap-5'>
      <h1 className='primary-title'>Dashboard</h1>
      <div>
        <p>メールアドレス:{email}</p>
      </div>
      <SkeltonButton href='/'>ホームへ</SkeltonButton>
    </div>
  );
}
