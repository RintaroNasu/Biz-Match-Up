import { UserProfileUpdate } from '../types';

export const editUserProfile = async (
  userId: number,
  body: UserProfileUpdate,
) => {
  const url = `http://localhost:3000/api/profile/${userId}`;

  const data = await fetch(url, {
    method: 'POST',
    body: JSON.stringify(body),
    headers: {
      'Content-Type': 'application/json',
    },
  });

  return data.json();
};
