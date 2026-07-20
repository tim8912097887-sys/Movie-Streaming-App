const BUTTON_SIZES = {
  sm: "px-3 py-1.5 text-sm rounded",
  md: "px-4 py-2 text-base rounded-md",
  lg: "px-6 py-3 text-lg rounded-lg",
};

const BUTTON_COLORS = {
  primary: "bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500",
  secondary: "bg-gray-200 text-gray-800 hover:bg-gray-300 focus:ring-gray-400",
  quaternary: "bg-gray-100 text-gray-800 hover:bg-gray-200 focus:ring-gray-400",
  success:
    "bg-emerald-500 text-white hover:bg-emerald-600 focus:ring-emerald-500",
  danger: "bg-rose-500 text-white hover:bg-rose-600 focus:ring-rose-500",
  warning: "bg-yellow-500 text-white hover:bg-yellow-600 focus:ring-yellow-500",
  info: "bg-sky-500 text-white hover:bg-sky-600 focus:ring-sky-500",
};

const BUTTON_TYPES = {
  normal:
    "font-semibold shadow-sm transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2",
  link: "bg-transparent text-blue-600 hover:underline shadow-none focus:ring-0",
};

type ButtonProps = {
  size: keyof typeof BUTTON_SIZES;
  color: keyof typeof BUTTON_COLORS;
  btnType: keyof typeof BUTTON_TYPES;
  buttonProps?: React.ButtonHTMLAttributes<HTMLButtonElement>;
  children: React.ReactNode;
};

const Button = ({
  size,
  color,
  btnType,
  buttonProps,
  children,
}: ButtonProps) => {
  return (
    <button
      {...buttonProps}
      className={`flex justify-center items-center cursor-pointer ${BUTTON_SIZES[size]} ${BUTTON_COLORS[color]} ${BUTTON_TYPES[btnType]}`}
    >
      {children}
    </button>
  );
};

export default Button;
