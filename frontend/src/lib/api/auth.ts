import { LoginUser, RegisterUser } from '../types';

export const signUp = async (body: RegisterUser) => {
  const url = 'http://localhost:3000/api/user/signup';

  const data = await fetch(url, {
    method: 'POST',
    body: JSON.stringify(body),
    headers: {
      'Content-Type': 'application/json',
    },
  });

  return data.json();
};

export const signIn = async (body: LoginUser) => {
  const url = 'http://localhost:3000/api/user/signin';

  const data = await fetch(url, {
    method: 'POST',
    body: JSON.stringify(body),
    headers: {
      'Content-Type': 'application/json',
    },
  });

  return data.json();
};
