import { UserProfileUpdate } from '../types';

export const editUserProfile = async (
  userId: number,
  body: UserProfileUpdate,
) => {
  const url = `http://localhost:8080/api/profile/${userId}`;

  const data = await fetch(url, {
    method: 'POST',
    body: JSON.stringify(body),
    headers: {
      'Content-Type': 'application/json',
    },
  });

  return data.json();
};

export const getUserProfile = async (userId: number) => {
  const url = `http://localhost:8080/api/profile/${userId}`;

  const data = await fetch(url, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  return data.json();
};
