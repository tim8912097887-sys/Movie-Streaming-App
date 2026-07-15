import FlexContainer from "../../ui/FlexContainer";
import PageLink from "../../ui/PageLink";

const LeftHeader = () => {
  return (
    <FlexContainer customClass="gap-8">
      <h1 className="text-4xl font-bold">Movie Streaming</h1>
      <FlexContainer customClass="gap-4">
        <PageLink target="/" text="Home" />
        <PageLink target="/recommendations" text="Recommendations" />
      </FlexContainer>
    </FlexContainer>
  );
};

export default LeftHeader;
