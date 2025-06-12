import { MatchItem } from '../types';

export const companyScrape = async (body: { companyUrl: string }) => {
  const url = 'http://localhost:8080/api/company-scrape';
  const token = localStorage.getItem('access_token');

  const data = await fetch(url, {
    method: 'POST',
    body: JSON.stringify(body),
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
  });

  return data.json();
};

export const generateReasons = async (body: {
  matchResult: MatchItem[];
  questions: {
    reasonInterest: string;
    attractiveService: string;
    relatedExperience: string;
  };
}) => {
  const url = 'http://localhost:8080/api/generate-reasons';
  const token = localStorage.getItem('access_token');

  const data = await fetch(url, {
    method: 'POST',
    body: JSON.stringify(body),
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
  });

  return data.json();
};

export const postReasons = async ({
  content,
  companyName,
  companyUrl,
}: {
  content: string;
  companyName: string;
  companyUrl: string;
}) => {
  const url = 'http://localhost:8080/api/reason';
  const token = localStorage.getItem('access_token');

  const data = await fetch(url, {
    method: 'POST',
    body: JSON.stringify({ content, companyName, companyUrl }),
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
  });

  return data.json();
};

export const getReasons = async () => {
  const token = localStorage.getItem('access_token');
  const res = await fetch('http://localhost:8080/api/reasons', {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return res.json();
};
