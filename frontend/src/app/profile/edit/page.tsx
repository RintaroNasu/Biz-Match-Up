'use client';
import Link from 'next/link';
import { ArrowLeft } from 'lucide-react';
import { PrimaryButton } from '../../../components/buttons/PrimaryButton';
import { useEffect, useState } from 'react';
import { jwtDecode } from 'jwt-decode';
import { useRouter } from 'next/navigation';
import { editUserProfile, getUserProfile } from '../../../lib/api/user';
import { UserProfileUpdate } from '../../../lib/types';
import { errorToast, successToast } from '@/lib/toast';

export default function EditProfile() {
  const router = useRouter();
  const [userId, setUserId] = useState<number | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      const token = localStorage.getItem('access_token');
      if (token && token.split('.').length === 3) {
        const decodedUser = jwtDecode<{
          id: number;
          email: string;
          iat: number;
          exp: number;
        }>(token);
        setUserId(decodedUser.id);
        if (decodedUser.id) {
          const res = await getUserProfile(decodedUser.id);
          if (res.user) {
            setForm({
              name: res.user.Name || '',
              desiredJobType: res.user.DesiredJobType || '',
              desiredLocation: res.user.DesiredLocation || '',
              desiredCompanySize: res.user.DesiredCompanySize || '',
              careerAxis1: res.user.CareerAxis1 || '',
              careerAxis2: res.user.CareerAxis2 || '',
              selfPr: res.user.SelfPr || '',
            });
          }
        }
      } else if (!token) {
        router.push('/');
      }
    };
    fetchData();
  }, [router]);

  const [form, setForm] = useState<UserProfileUpdate>({
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

  const onSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!userId) return;
    try {
      const res = await editUserProfile(userId, form);
      if (res.success) {
        successToast('プロフィールを更新しました');
      } else {
        errorToast('プロフィールの更新に失敗しました');
      }
    } catch (err) {
      console.error('プロフィール更新エラー:', err);
    }
  };

  return (
    <div className='py-8 px-4'>
      <Link
        href='/dashboard'
        className='flex items-center gap-2 text-blue-600 mb-6'
      >
        <ArrowLeft />
        <span>ダッシュボードに戻る</span>
      </Link>

      <h1 className='text-3xl font-bold mb-8'>プロフィール編集</h1>
      <div className='px-32'>
        <div className='max-w-2xl mx-auto bg-white shadow-md rounded-lg p-6'>
          <h2 className='text-xl font-semibold mb-4'>個人情報</h2>

          <form onSubmit={onSubmit} className='space-y-6'>
            <div>
              <label className='block text-sm font-medium text-gray-700'>
                氏名
              </label>
              <input
                name='name'
                value={form.name}
                onChange={onChangeForm}
                placeholder='山田 太郎'
                className='mt-1 block w-full border border-gray-300 rounded-md p-2'
              />
            </div>

            <div>
              <label className='block text-sm font-medium text-gray-700'>
                志望職種
              </label>
              <input
                name='desiredJobType'
                value={form.desiredJobType}
                onChange={onChangeForm}
                type='text'
                placeholder='エンジニア / マーケターなど'
                className='mt-1 block w-full border border-gray-300 rounded-md p-2'
              />
            </div>

            <div>
              <label className='block text-sm font-medium text-gray-700'>
                志望勤務地
              </label>
              <input
                name='desiredLocation'
                value={form.desiredLocation}
                onChange={onChangeForm}
                placeholder='東京 / 福岡など'
                className='mt-1 block w-full border border-gray-300 rounded-md p-2'
              />
            </div>

            <div>
              <label className='block text-sm font-medium text-gray-700'>
                志望企業の規模
              </label>
              <input
                name='desiredCompanySize'
                value={form.desiredCompanySize}
                onChange={onChangeForm}
                placeholder='大手 / ベンチャーなど'
                className='mt-1 block w-full border border-gray-300 rounded-md p-2'
              />
            </div>

            <div>
              <label className='block text-sm font-medium text-gray-700'>
                就活軸の1つ目
              </label>
              <textarea
                name='careerAxis1'
                value={form.careerAxis1}
                onChange={onChangeForm}
                placeholder='例：社会貢献性が高い仕事がしたい'
                className='mt-1 block w-full border border-gray-300 rounded-md p-2'
              />
            </div>

            <div>
              <label className='block text-sm font-medium text-gray-700'>
                就活軸の2つ目
              </label>
              <textarea
                name='careerAxis2'
                value={form.careerAxis2}
                onChange={onChangeForm}
                placeholder='例：成長できる環境が良い'
                className='mt-1 block w-full border border-gray-300 rounded-md p-2'
              />
            </div>

            <div>
              <label className='block text-sm font-medium text-gray-700'>
                自己PR
              </label>
              <textarea
                name='selfPr'
                value={form.selfPr}
                onChange={onChangeForm}
                placeholder='自己PRを入力してください'
                className='mt-1 block w-full border border-gray-300 rounded-md p-2 min-h-[150px]'
              />
            </div>

            <div className='text-right'>
              <PrimaryButton>保存する</PrimaryButton>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
