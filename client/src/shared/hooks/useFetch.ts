import { useCallback, useState } from "react";
import z from "zod";

type FormStatus<T> = {
  isFetching: boolean;
  isSuccess: boolean;
  error: string | null;
  fetchedData: T | null;
};

type Props<T, U> = {
  fetchFunction: (data: T) => Promise<U>;
  validateSchema?: z.ZodObject<any>;
};

const useFetch = <T, U>({ fetchFunction, validateSchema }: Props<T, U>) => {
  const [status, setStatus] = useState<FormStatus<U>>({
    isFetching: false,
    isSuccess: false,
    error: null,
    fetchedData: null,
  });

  const handleFetch = useCallback(
    async (data: T) => {
      setStatus((preState) => ({
        ...preState,
        isFetching: true,
        isSuccess: false,
        error: null,
      }));
      try {
        // Validate form data before submission if needed
        if (validateSchema) {
          const result = validateSchema.safeParse(data);
          console.log("validation result", result);
          if (!result.success) {
            setStatus({
              isFetching: false,
              isSuccess: false,
              error: result.error.issues[0].message,
              fetchedData: null,
            });
            return;
          }
        }
        const response = await fetchFunction(data);
        setStatus({
          isFetching: false,
          isSuccess: true,
          error: null,
          fetchedData: response,
        });
      } catch (error: any) {
        console.error(`Error submitting form: ${error.message}`, error);
        setStatus({
          isFetching: false,
          isSuccess: false,
          error: error.message || "An error occurred during form submission",
          fetchedData: null,
        });
      }
    },
    [fetchFunction, validateSchema],
  );

  return {
    status,
    handleFetch,
  };
};

export default useFetch;
