import React from 'react';
import { Outlet } from 'react-router-dom';
import { Sections_teacher } from './sections_teacher';

export function MainLayout_teacher() {
  return (
    <div className="app-layout">
      <aside className="sidebar">
        <Sections_teacher />
      </aside>
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  );
}
