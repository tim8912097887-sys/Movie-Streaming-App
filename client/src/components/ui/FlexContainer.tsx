type FlexContainerProps = {
  customClass: string;
  children: React.ReactNode;
};

const FlexContainer = ({ customClass, children }: FlexContainerProps) => {
  return <div className={`flex items-center ${customClass}`}>{children}</div>;
};

export default FlexContainer;
