import FlexContainer from "../../ui/FlexContainer";
import LeftHeader from "./LeftHeader";
import RightHeader from "./RightHeader";

const Header = () => {
  return (
    <FlexContainer customClass="bg-black text-white justify-between px-4 sm:px-8 md:px-12 lg:px-16 py-6">
      <LeftHeader />
      <RightHeader />
    </FlexContainer>
  );
};

export default Header;
