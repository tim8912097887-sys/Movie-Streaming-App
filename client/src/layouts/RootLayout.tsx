import Header from "../components/layout/Header";

const RootLayout = () => {
  return (
    <div className="flex flex-col min-h-screen bg-white text-black font-sans antialiased">
      <Header />
      <main className="flex-1 px-4 py-6 md:px-8 max-w-7xl mx-auto w-full"></main>
    </div>
  );
};

export default RootLayout;
