type InputGroupProps = {
  name: string;
  label: string;
  children: React.ReactNode;
};
const InputGroup = ({ name, label, children }: InputGroupProps) => {
  return (
    <div className="flex flex-col gap-3">
      <label htmlFor={name}>{label}:</label>
      {children}
    </div>
  );
};

export default InputGroup;
