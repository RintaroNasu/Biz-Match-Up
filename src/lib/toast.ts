import { toast } from 'react-toastify';

export const errorToast = (message: string) => {
  return toast.error(message, {
    autoClose: 1500,
    closeButton: false,
    hideProgressBar: true,
    pauseOnHover: false,
  });
};

export const successToast = (message: string) => {
  return toast.success(message, {
    autoClose: 1500,
    closeButton: false,
    hideProgressBar: true,
    pauseOnHover: false,
  });
};
