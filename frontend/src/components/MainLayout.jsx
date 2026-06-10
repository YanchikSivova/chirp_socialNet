
import "../styles/mainLayout.css";
import NavigationMenu from "./NavigationMenu.jsx";
function MainLayout({ children }) {
    return (
        <main className="main-page">
            <header className="main-header">
                <div
                    className="logo-text"
                    aria-label="Чирик"
                >
                    ЧИРИК
                </div>
            </header>

            <section className="main-content">
                {children}
            </section>

            <NavigationMenu />
        </main>
    );
}

export default MainLayout;