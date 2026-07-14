import React from "react";

// Define the available prop types
type SpinnerSize = "xs" | "sm" | "md" | "lg" | "xl";
type SpinnerColor = "primary" | "secondary" | "success" | "danger" | "white";

interface SpinnerProps {
  size?: SpinnerSize;
  color?: SpinnerColor;
  className?: string; // For any extra custom overrides
}

// Predefined Tailwind mappings
const sizeMap: Record<SpinnerSize, string> = {
  xs: "h-4 w-4 border-2",
  sm: "h-6 w-6 border-2",
  md: "h-10 w-10 border-3",
  lg: "h-16 w-16 border-4",
  xl: "h-24 w-24 border-8",
};

const colorMap: Record<SpinnerColor, string> = {
  primary: "border-blue-600 border-t-transparent",
  secondary: "border-gray-600 border-t-transparent",
  success: "border-emerald-500 border-t-transparent",
  danger: "border-rose-500 border-t-transparent",
  white: "border-white border-t-transparent",
};

export const Spinner: React.FC<SpinnerProps> = ({
  size = "md",
  color = "primary",
  className = "",
}) => {
  return (
    <div
      className={`
        animate-spin 
        rounded-full 
        ${sizeMap[size]} 
        ${colorMap[color]} 
        ${className}
      `}
      role="status"
      aria-label="loading"
    >
      <span className="sr-only">Loading...</span>
    </div>
  );
};
