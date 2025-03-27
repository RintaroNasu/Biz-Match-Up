export const companyScrape = async (body: { companyUrl: string }) => {
  const url = 'http://localhost:3000/api/companyScrape';

  const data = await fetch(url, {
    method: 'POST',
    body: JSON.stringify(body),
    headers: {
      'Content-Type': 'application/json',
    },
  });

  return data.json();
};
