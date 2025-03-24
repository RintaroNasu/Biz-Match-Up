export const signUp = async (body: { email: string; password: string }) => {
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

export const signIn = async (body: { email: string; password: string }) => {
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
