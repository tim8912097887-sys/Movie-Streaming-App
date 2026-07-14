import type { PropsWithChildren } from "react";

type FormProps = {
  children: React.ReactNode;
  onSubmit?: (e: React.SubmitEvent<HTMLFormElement>) => void;
};

const Form = ({ children, onSubmit }: FormProps) => {
  return (
    <form className="w-full space-y-10" onSubmit={onSubmit}>
      {children}
    </form>
  );
};

const Head = ({ children }: PropsWithChildren) => {
  return <div className="text-2xl font-bold text-center">{children}</div>;
};

const Content = ({ children }: PropsWithChildren) => {
  return <div className="flex flex-col gap-4">{children}</div>;
};

const Footer = ({ children }: PropsWithChildren) => {
  return <div className="flex flex-col gap-2">{children}</div>;
};

Form.Content = Content;
Form.Head = Head;
Form.Footer = Footer;

export default Form;
