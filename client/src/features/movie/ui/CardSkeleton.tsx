import Skeleton from "react-loading-skeleton";
import Card from "./Card";

const CardSkeleton = () => {
  return (
    <Card>
      <Card.Head>
        <div className="h-full w-full">
          <Skeleton
            height="100%"
            width="100%"
            containerClassName="block w-full h-full"
          />
        </div>
      </Card.Head>
      <Card.Body>
        <Skeleton height={28} width="70%" />
        <Skeleton count={3} />
      </Card.Body>
      <Card.Footer>
        <Skeleton width={80} />
      </Card.Footer>
    </Card>
  );
};

export default CardSkeleton;
