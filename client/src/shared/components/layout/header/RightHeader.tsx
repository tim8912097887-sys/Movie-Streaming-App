import FlexContainer from "../../ui/FlexContainer";
import PageLink from "../../ui/PageLink";

const RightHeader = () => {
  return (
    <FlexContainer customClass="gap-4">
      <PageLink target="/login" text="Login" />
      <PageLink target="/register" text="Register" />
    </FlexContainer>
  );
};

export default RightHeader;
