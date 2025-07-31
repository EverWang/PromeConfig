import React, { useState, useEffect } from 'react';
import { AlertTriangle, Database, Settings, Zap, RefreshCw, Eye, Download, Upload, Cpu } from 'lucide-react';
import { AuthWrapper } from './components/AuthWrapper';
import { Sidebar } from './components/Sidebar';
import { Dashboard } from './components/Dashboard';
import { TargetManagement } from './components/TargetManagement';
import { AlertRuleManagement } from './components/AlertRuleManagement';
import { ConfigPreview } from './components/ConfigPreview';
import { PrometheusAPI } from './components/PrometheusAPI';
import { targetService } from './services/targetService';
import { alertRuleService } from './services/alertRuleService';
import authService from './services/authService';
import type { Target, AlertRule } from './types';

type ActiveView = 'dashboard' | 'targets' | 'alerts' | 'preview' | 'api';

function App() {
  const [activeView, setActiveView] = useState<ActiveView>('dashboard');
  const [targets, setTargets] = useState<Target[]>([]);
  const [alertRules, setAlertRules] = useState<AlertRule[]>([]);
  const [loading, setLoading] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  // 监听认证状态变化
  useEffect(() => {
    const checkAuth = (shouldLoadData = false) => {
      const authenticated = authService.isAuthenticated();
      const wasAuthenticated = isAuthenticated;
      setIsAuthenticated(authenticated);
      
      // 只有在认证状态发生变化或明确要求加载数据时才重新加载
      if (authenticated && (shouldLoadData || !wasAuthenticated)) {
        loadData();
      } else if (!authenticated && wasAuthenticated) {
        setLoading(false);
        setTargets([]);
        setAlertRules([]);
      }
    };

    // 初始检查并加载数据
    checkAuth(true);

    // 监听存储变化（当用户在其他标签页登录/登出时）
    const handleStorageChange = () => {
      checkAuth(true);
    };

    window.addEventListener('storage', handleStorageChange);
    
    // 定期检查认证状态，但不重新加载数据
    const interval = setInterval(() => checkAuth(false), 5000);

    return () => {
      window.removeEventListener('storage', handleStorageChange);
      clearInterval(interval);
    };
  }, [isAuthenticated]);

  const loadData = async () => {
    try {
      setLoading(true);
      const [targetsData, alertRulesData] = await Promise.all([
        targetService.getTargets(),
        alertRuleService.getAlertRules()
      ]);
      setTargets(targetsData);
      setAlertRules(alertRulesData);
    } catch (error) {
      console.error('Error loading data:', error);
    } finally {
      setLoading(false);
    }
  };

  const menuItems = [
    { id: 'dashboard', label: 'Dashboard', icon: Cpu },
    { id: 'targets', label: 'Targets', icon: Database },
    { id: 'alerts', label: 'Alert Rules', icon: AlertTriangle },
    { id: 'preview', label: 'Config Preview', icon: Eye },
    { id: 'api', label: 'API Management', icon: RefreshCw },
  ];

  const renderActiveView = () => {
    if (loading) {
      return (
        <div className="flex items-center justify-center h-full">
          <div className="text-white">Loading...</div>
        </div>
      );
    }

    switch (activeView) {
      case 'dashboard':
        return <Dashboard targets={targets} alertRules={alertRules} />;
      case 'targets':
        return <TargetManagement targets={targets} onDataChange={loadData} />;
      case 'alerts':
        return <AlertRuleManagement alertRules={alertRules} onDataChange={loadData} />;
      case 'preview':
        return <ConfigPreview targets={targets} alertRules={alertRules} />;
      case 'api':
        return <PrometheusAPI targets={targets} alertRules={alertRules} />;
      default:
        return <Dashboard targets={targets} alertRules={alertRules} />;
    }
  };

  return (
    <AuthWrapper>
      <div className="flex h-screen bg-gray-900 text-gray-100">
        <Sidebar
          menuItems={menuItems}
          activeView={activeView}
          setActiveView={setActiveView}
        />
        <main className="flex-1 overflow-auto">
          {renderActiveView()}
        </main>
      </div>
    </AuthWrapper>
  );
}

export default App;