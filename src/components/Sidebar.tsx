import React from 'react';
import { Settings, Zap } from 'lucide-react';

interface MenuItem {
  id: string;
  label: string;
  icon: React.ComponentType<any>;
}

interface SidebarProps {
  menuItems: MenuItem[];
  activeView: string;
  setActiveView: (view: string) => void;
}

export const Sidebar: React.FC<SidebarProps> = ({ menuItems, activeView, setActiveView }) => {
  return (
    <div className="w-64 bg-gray-800 border-r border-gray-700 flex flex-col">
      <div className="p-6 border-b border-gray-700">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
            <Zap className="w-5 h-5 text-white" />
          </div>
          <div>
            <h1 className="text-xl font-bold text-white">PromeConfig</h1>
            <p className="text-sm text-gray-400">Configuration Manager</p>
          </div>
        </div>
      </div>

      <nav className="flex-1 p-4">
        <ul className="space-y-2">
          {menuItems.map((item) => {
            const Icon = item.icon;
            const isActive = activeView === item.id;
            
            return (
              <li key={item.id}>
                <button
                  onClick={() => setActiveView(item.id)}
                  className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-all duration-200 ${
                    isActive
                      ? 'bg-blue-600 text-white shadow-lg'
                      : 'text-gray-300 hover:bg-gray-700 hover:text-white'
                  }`}
                >
                  <Icon className="w-5 h-5" />
                  <span className="font-medium">{item.label}</span>
                </button>
              </li>
            );
          })}
        </ul>
      </nav>

      <div className="p-4 border-t border-gray-700">
        <div className="flex items-center gap-3 p-3 bg-gray-700 rounded-lg">
          <Settings className="w-5 h-5 text-gray-400" />
          <div className="flex-1">
            <p className="text-sm font-medium text-white">Prometheus</p>
            <p className="text-xs text-gray-400">Connected</p>
          </div>
          <div className="w-2 h-2 bg-green-500 rounded-full"></div>
        </div>
      </div>
    </div>
  );
};