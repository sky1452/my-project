import { Outlet } from 'react-router-dom';
import { SectionsTeacher } from './sections_teacher';

export function MainLayoutTeacher() {
  return (
    <div className="app-layout">
      <aside className="sidebar">
        <SectionsTeacher />
      </aside>
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  );
}