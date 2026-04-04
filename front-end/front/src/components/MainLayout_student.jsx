import { Outlet } from 'react-router-dom';
import { SectionsStudent } from './sections_student';

export function MainLayoutStudent() {
  return (
    <div className="app-layout">
      <aside className="sidebar">
        <SectionsStudent />
      </aside>
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  );
}