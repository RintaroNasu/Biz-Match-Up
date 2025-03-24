import Link from 'next/link';

type Props = {
  children: React.ReactNode;
  href: string;
};

export const SkeltonButton = (props: Props) => {
  return (
    <Link
      href={props.href}
      className='px-6 py-3 text-blue-600 border border-blue-600 hover:bg-blue-100 rounded-lg text-lg font-semibold shadow-md transition cursor-pointer'
    >
      {props.children}
    </Link>
  );
};
