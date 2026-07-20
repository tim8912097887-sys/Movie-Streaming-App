import { useEffect } from "react";
import { logoutUser } from "../../../../features/auth/api/request";
import { useAuthStore } from "../../../../features/auth/store/store";
import useFetch from "../../../hooks/useFetch";
import Button from "../../ui/Button";
import FlexContainer from "../../ui/FlexContainer";
import PageLink from "../../ui/PageLink";
import { toast } from "react-toastify";

const RightHeader = () => {
  const token = useAuthStore((state) => state.token);
  const logout = useAuthStore((state) => state.logout);
  const { handleFetch, status } = useFetch({
    fetchFunction: logoutUser,
  });

  const handleLogout = async () => {
    await handleFetch(token!);
  };

  useEffect(() => {
    if (status.error) {
      toast.success("Logout failed", {
        autoClose: 1500,
        position: "top-right",
      });
    }

    if (status.isSuccess) {
      toast.success("Logout successfully", {
        autoClose: 1500,
        position: "top-right",
        onClose: () => logout(),
      });
    }
  }, [status.error, status.isSuccess]);

  return (
    <FlexContainer customClass="gap-4">
      {token ? (
        <>
          <Button
            size="sm"
            color="info"
            btnType="link"
            buttonProps={{
              onClick: () => handleLogout(),
              disabled: status.isFetching,
            }}
            children="Logout"
          />
        </>
      ) : (
        <>
          <PageLink target="/login" text="Login" />
          <PageLink target="/register" text="Register" />
        </>
      )}
    </FlexContainer>
  );
};

export default RightHeader;
