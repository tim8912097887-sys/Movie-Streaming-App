const BUTTON_SIZES = {
  sm: "px-3 py-1.5 text-sm rounded",
  md: "px-4 py-2 text-base rounded-md",
  lg: "px-6 py-3 text-lg rounded-lg",
};

const BUTTON_COLORS = {
  primary: "bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500",
  secondary: "bg-gray-200 text-gray-800 hover:bg-gray-300 focus:ring-gray-400",
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
      className={`cursor-pointer ${BUTTON_SIZES[size]} ${BUTTON_COLORS[color]} ${BUTTON_TYPES[btnType]}`}
    >
      {children}
    </button>
  );
};

export default Button;
