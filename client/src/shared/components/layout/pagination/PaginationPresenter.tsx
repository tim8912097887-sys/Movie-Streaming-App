import PageButton from "../../ui/PageButton";
import PaginationSkeleton from "./PaginationSkeleton";

type PaginationPresenterProps = {
  totalPage: number;
  currentPage: number;
  onPageChange: (page: number) => void;
  isLoading: boolean;
  setAttemptedPage: (page: number) => void;
};

const PaginationPresenter = ({
  totalPage,
  currentPage,
  onPageChange,
  isLoading,
  setAttemptedPage,
}: PaginationPresenterProps) => {
  if (isLoading && (totalPage === 0 || currentPage === 0))
    return <PaginationSkeleton />;

  if (totalPage === 0 || currentPage === 0)
    return (
      <div className="text-gray-400 text-center p-4">No pages to display.</div>
    );

  const pageArray = [currentPage - 1, currentPage, currentPage + 1];

  return (
    <div className="flex overflow-hidden border border-gray-300 shadow-2xs rounded">
      {currentPage == 1 ? (
        <PageButton customClass="w-15" buttonProps={{ disabled: true }}>
          Prev
        </PageButton>
      ) : (
        <PageButton
          customClass="w-15"
          buttonProps={{
            disabled: isLoading,
            onClick: () => {
              onPageChange(currentPage - 1);
              setAttemptedPage(currentPage - 1);
            },
          }}
        >
          Prev
        </PageButton>
      )}
      <PageButton
        buttonProps={{
          disabled: isLoading,
          onClick: () => {
            onPageChange(1);
            setAttemptedPage(1);
          },
        }}
        isActive={currentPage == 1}
      >
        1
      </PageButton>

      {currentPage - 2 == 2 && (
        <PageButton
          buttonProps={{
            disabled: isLoading,
            onClick: () => {
              onPageChange(2);
              setAttemptedPage(2);
            },
          }}
        >
          2
        </PageButton>
      )}

      {currentPage > 4 && (
        <PageButton buttonProps={{ disabled: true }}>...</PageButton>
      )}

      {pageArray.map((page) =>
        page > 1 && page < totalPage ? (
          <PageButton
            buttonProps={{
              disabled: isLoading,
              onClick: () => {
                onPageChange(page);
                setAttemptedPage(page);
              },
            }}
            isActive={page == currentPage}
            key={page}
          >
            {page}
          </PageButton>
        ) : null,
      )}

      {currentPage + 2 == totalPage - 1 && (
        <PageButton
          buttonProps={{
            disabled: isLoading,
            onClick: () => {
              onPageChange(totalPage - 1);
              setAttemptedPage(totalPage - 1);
            },
          }}
        >
          {totalPage - 1}
        </PageButton>
      )}

      {currentPage < totalPage - 3 && totalPage > 5 && (
        <PageButton buttonProps={{ disabled: true }}>...</PageButton>
      )}

      {totalPage > 1 && (
        <PageButton
          isActive={currentPage == totalPage}
          buttonProps={{
            disabled: isLoading,
            onClick: () => {
              onPageChange(totalPage);
              setAttemptedPage(totalPage);
            },
          }}
        >
          {totalPage}
        </PageButton>
      )}
      {currentPage == totalPage ? (
        <PageButton customClass="w-15" buttonProps={{ disabled: true }}>
          Next
        </PageButton>
      ) : (
        <PageButton
          customClass="w-15"
          buttonProps={{
            disabled: isLoading,
            onClick: () => {
              onPageChange(currentPage + 1);
              setAttemptedPage(currentPage + 1);
            },
          }}
        >
          Next
        </PageButton>
      )}
    </div>
  );
};

export default PaginationPresenter;
