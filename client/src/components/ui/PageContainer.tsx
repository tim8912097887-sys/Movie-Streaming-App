type PageContainerProps = {
  children: React.ReactNode;
  customClass?: string;
};

const PageContainer = ({ children, customClass }: PageContainerProps) => {
  return (
    <div
      className={`flex min-h-full items-center justify-center py-8 md:py-12 ${customClass}`}
    >
      {children}
    </div>
  );
};

export default PageContainer;
