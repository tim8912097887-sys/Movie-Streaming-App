type ErrorTextProps = {
  children: React.ReactNode;
};

const ErrorText = ({ children }: ErrorTextProps) => {
  return <p className="text-sm text-red-600">{children}</p>;
};

export default ErrorText;
