import { twMerge } from "tailwind-merge";

type PageContainerProps = {
  children: React.ReactNode;
  customClass?: string;
};

const PageContainer = ({ children, customClass }: PageContainerProps) => {
  return (
    <div
      className={twMerge(
        `flex min-h-96 items-center justify-center py-8 md:py-12`,
        customClass,
      )}
    >
      {children}
    </div>
  );
};

export default PageContainer;
