import { NavLink } from "react-router";

type PageLinkProps = {
  target: string;
  text: string;
};

const PageLink = ({ target, text }: PageLinkProps) => {
  return (
    <NavLink
      to={target}
      className={({ isActive }) =>
        isActive ? "text-white text-xl" : "text-gray-400 text-xl"
      }
    >
      {text}
    </NavLink>
  );
};

export default PageLink;
