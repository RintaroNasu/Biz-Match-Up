import { LoginUser, RegisterUser } from '../types';

export const signUp = async (body: RegisterUser) => {
  const url = 'http://localhost:8080/api/signup';

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
  const url = 'http://localhost:8080/api/signin';

  const data = await fetch(url, {
    method: 'POST',
    body: JSON.stringify(body),
    headers: {
      'Content-Type': 'application/json',
    },
  });

  return data.json();
};
