import { useState } from "react";
import z from "zod";

type UseFormDataProps<T extends Record<string, unknown>> = {
  initialValues: T;
  schemaValidater: z.ZodObject<any>;
};

export const useFormData = <T extends Record<string, unknown>>({
  initialValues,
  schemaValidater,
}: UseFormDataProps<T>) => {
  const [formData, setFormData] = useState<T>(initialValues);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;

    const currentValue = {
      ...formData,
      [name]: value,
    };
    setFormData(currentValue);
    const result = schemaValidater.safeParse(currentValue);
    if (!result.success) {
      const validateError = z.flattenError(result.error).fieldErrors;
      const errorMessage = validateError[name]?.[0];
      if (errorMessage) {
        setErrors((preError) => ({
          ...preError,
          [name]: errorMessage,
        }));
      } else {
        setErrors((preError) => {
          const newError = { ...preError };
          delete newError[name];
          return newError;
        });
      }
    } else {
      // Clear error when input is valid
      setErrors({});
    }
    setFormData({ ...formData, [name]: value });
  };

  return { formData, handleChange, errors };
};
