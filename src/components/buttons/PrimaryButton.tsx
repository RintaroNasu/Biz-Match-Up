type Props = {
  className?: string;
  children: React.ReactNode;
  disabled?: boolean;
  onClick?: () => void;
};

export const PrimaryButton = (props: Props) => {
  return (
    <button
      onClick={props.onClick}
      disabled={props.disabled}
      className={`px-6 py-3 text-white  rounded-lg text-lg font-semibold shadow-md transition cursor-pointer
    ${props.disabled ? 'cursor-not-allowed bg-blue-200' : 'bg-blue-500 hover:bg-blue-400'} 
    ${props.className}`}
    >
      {props.children}
    </button>
  );
};
