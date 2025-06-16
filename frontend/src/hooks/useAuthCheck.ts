'use client';

import { useEffect, useState } from 'react';
import { jwtDecode } from 'jwt-decode';
import { useRouter } from 'next/navigation';
import { error } from 'console';
import { errorToast } from '@/lib/toast';

type DecodedUser = {
  id: number;
  email: string;
  iat: number;
  exp: number;
};

export const useAuthCheck = () => {
  const [user, setUser] = useState<DecodedUser | null>(null);
  const router = useRouter();

  useEffect(() => {
    const token = localStorage.getItem('access_token');
    if (token && token.split('.').length === 3) {
      const decoded = jwtDecode<DecodedUser>(token);
      const now = Date.now() / 1000;
      if (decoded.exp < now) {
        localStorage.removeItem('access_token');
        errorToast('セッションが切れました。再度ログインしてください。');
        router.push('/');
      } else {
        setUser(decoded);
      }
    } else {
      router.push('/');
    }
  }, [router]);

  return { user };
};
