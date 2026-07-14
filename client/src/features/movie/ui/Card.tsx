import type { PropsWithChildren } from "react";

const Card = ({ children }: PropsWithChildren) => {
  return (
    <div className="w-full max-w-64 aspect-2/3 bg-black rounded-xl shadow-lg relative overflow-hidden flex flex-col justify-end hover:scale-105 transition-all duration-300">
      {children}
    </div>
  );
};

const Head = ({ children }: PropsWithChildren) => {
  return <div className="absolute inset-0 w-full h-full z-0">{children}</div>;
};

const Body = ({ children }: PropsWithChildren) => {
  return (
    <div className="w-full p-4 flex flex-col gap-2 z-10 bg-linear-to-t from-black via-black/80 to-transparent pt-16">
      {children}
    </div>
  );
};

const Footer = ({ children }: PropsWithChildren) => {
  return (
    <div className="w-full px-4 pb-4 flex gap-2 z-10 bg-black/80">
      {children}
    </div>
  );
};

Card.Head = Head;
Card.Body = Body;
Card.Footer = Footer;

export default Card;
