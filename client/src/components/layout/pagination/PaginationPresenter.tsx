import PageButton from "../../ui/PageButton";

type PaginationPresenterProps = {
  totalPage: number;
  currentPage: number;
  onPageChange: (page: number) => void;
};

const PaginationPresenter = ({
  totalPage,
  currentPage,
  onPageChange,
}: PaginationPresenterProps) => {
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
          buttonProps={{ onClick: () => onPageChange(currentPage - 1) }}
        >
          Prev
        </PageButton>
      )}
      <PageButton
        buttonProps={{ onClick: () => onPageChange(1) }}
        isActive={currentPage == 1}
      >
        1
      </PageButton>

      {currentPage - 2 == 2 && (
        <PageButton buttonProps={{ onClick: () => onPageChange(2) }}>
          2
        </PageButton>
      )}

      {currentPage > 4 && (
        <PageButton buttonProps={{ disabled: true }}>...</PageButton>
      )}

      {pageArray.map((page) =>
        page > 1 && page < totalPage ? (
          <PageButton
            buttonProps={{ onClick: () => onPageChange(page) }}
            isActive={page == currentPage}
            key={page}
          >
            {page}
          </PageButton>
        ) : null,
      )}

      {currentPage + 2 == totalPage - 1 && (
        <PageButton
          buttonProps={{ onClick: () => onPageChange(totalPage - 1) }}
        >
          {totalPage - 1}
        </PageButton>
      )}

      {currentPage < totalPage - 3 && totalPage > 5 && (
        <PageButton buttonProps={{ disabled: true }}>...</PageButton>
      )}

      <PageButton
        isActive={currentPage == totalPage}
        buttonProps={{ onClick: () => onPageChange(totalPage) }}
      >
        {totalPage}
      </PageButton>
      {currentPage == totalPage ? (
        <PageButton customClass="w-15" buttonProps={{ disabled: true }}>
          Next
        </PageButton>
      ) : (
        <PageButton
          customClass="w-15"
          buttonProps={{ onClick: () => onPageChange(currentPage + 1) }}
        >
          Next
        </PageButton>
      )}
    </div>
  );
};

export default PaginationPresenter;
