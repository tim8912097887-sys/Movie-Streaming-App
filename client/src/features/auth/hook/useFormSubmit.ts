import { useCallback, useState } from "react";
import z from "zod";

type FormStatus = {
  isSubmitting: boolean;
  isSuccess: boolean;
  error: string | null;
};

type Props<T extends Record<string, unknown>> = {
  submitFunction: (data: T) => Promise<void>;
  validateSchema?: z.ZodObject<any>;
};

const useFormSubmit = <T extends Record<string, unknown>>({
  submitFunction,
  validateSchema,
}: Props<T>) => {
  const [status, setStatus] = useState<FormStatus>({
    isSubmitting: false,
    isSuccess: false,
    error: null,
  });

  const handleSubmit = useCallback(
    async (formData: T) => {
      setStatus({
        isSubmitting: true,
        isSuccess: false,
        error: null,
      });
      try {
        // Validate form data before submission if needed
        if (validateSchema) {
          const result = validateSchema.safeParse(formData);
          console.log("validation result", result);
          if (!result.success) {
            setStatus({
              isSubmitting: false,
              isSuccess: false,
              error: result.error.issues[0].message,
            });
            return;
          }
        }
        await submitFunction(formData);
        setStatus({
          isSubmitting: false,
          isSuccess: true,
          error: null,
        });
      } catch (error: any) {
        console.error(`Error submitting form: ${error.message}`, error);
        setStatus({
          isSubmitting: false,
          isSuccess: false,
          error: error.message || "An error occurred during form submission",
        });
      }
    },
    [submitFunction, validateSchema],
  );

  return {
    status,
    handleSubmit,
  };
};

export default useFormSubmit;
