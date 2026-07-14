import { twMerge } from "tailwind-merge";

type PageButtonProps = {
  buttonProps?: React.ButtonHTMLAttributes<HTMLButtonElement>;
  children: React.ReactNode;
  isActive?: boolean;
  customClass?: string;
};

const PageButton = ({
  buttonProps,
  children,
  isActive = false,
  customClass,
}: PageButtonProps) => {
  return (
    <button
      {...buttonProps}
      disabled={buttonProps?.disabled || isActive}
      className={twMerge(
        `rounded-md ${isActive && "border"} bg-white flex justify-center items-center text-sm font-medium ${buttonProps?.disabled ? "text-gray-400" : "text-gray-700"} ${!buttonProps?.disabled && !isActive && "hover:bg-gray-500 cursor-pointer"} w-10 h-10`,
        customClass,
      )}
    >
      <span className={`${isActive ? "font-extrabold" : "font-medium"}`}>
        {children}
      </span>
    </button>
  );
};

export default PageButton;
