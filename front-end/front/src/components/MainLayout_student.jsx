import React from 'react';
import { Outlet } from 'react-router-dom';
import { Sections_student } from './sections_student';

export function MainLayout_student() {
  return (
    <div className="app-layout">
      <aside className="sidebar">
        <Sections_student />
      </aside>
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  );
}
