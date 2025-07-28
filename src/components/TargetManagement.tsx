import React, { useState } from 'react';
import { Plus, Edit, Trash2, Database, Settings, ChevronDown, ChevronRight } from 'lucide-react';
import { TargetService } from '../services/targetService';
import type { Target } from '../lib/supabase';

interface TargetManagementProps {
  targets: Target[];
  onDataChange: () => void;
}

export const TargetManagement: React.FC<TargetManagementProps> = ({ targets, onDataChange }) => {
  const [showForm, setShowForm] = useState(false);
  const [editingTarget, setEditingTarget] = useState<Target | null>(null);
  const [showRelabelConfig, setShowRelabelConfig] = useState(false);
  const [saving, setSaving] = useState(false);
  const [formData, setFormData] = useState({
    job_name: '',
    targets: '',
    scrape_interval: '15s',
    metrics_path: '/metrics',
    relabel_configs: [] as Array<{
      source_labels?: string[];
      separator?: string;
      target_label?: string;
      regex?: string;
      modulus?: number;
      replacement?: string;
      action?: 'replace' | 'keep' | 'drop' | 'hashmod' | 'labelmap' | 'labeldrop' | 'labelkeep';
    }>,
    metric_relabel_configs: [] as Array<{
      source_labels?: string[];
      separator?: string;
      target_label?: string;
      regex?: string;
      modulus?: number;
      replacement?: string;
      action?: 'replace' | 'keep' | 'drop' | 'hashmod' | 'labelmap' | 'labeldrop' | 'labelkeep';
    }>
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    
    const targetData = {
      job_name: formData.job_name,
      targets: formData.targets.split(',').map(t => t.trim()),
      scrape_interval: formData.scrape_interval,
      metrics_path: formData.metrics_path,
      relabel_configs: formData.relabel_configs.length > 0 ? formData.relabel_configs : undefined,
      metric_relabel_configs: formData.metric_relabel_configs.length > 0 ? formData.metric_relabel_configs : undefined,
    };

    try {
      if (editingTarget) {
        await TargetService.updateTarget(editingTarget.id, targetData);
      } else {
        await TargetService.createTarget(targetData);
      }
      onDataChange();
      resetForm();
    } catch (error) {
      console.error('Error saving target:', error);
      alert('保存失败，请重试');
    } finally {
      setSaving(false);
    }
  };

  const resetForm = () => {
    setFormData({
      job_name: '',
      targets: '',
      scrape_interval: '15s',
      metrics_path: '/metrics',
      relabel_configs: [],
      metric_relabel_configs: []
    });
    setShowForm(false);
    setEditingTarget(null);
    setShowRelabelConfig(false);
  };

  const handleEdit = (target: Target) => {
    setEditingTarget(target);
    setFormData({
      job_name: target.job_name,
      targets: target.targets.join(', '),
      scrape_interval: target.scrape_interval,
      metrics_path: target.metrics_path,
      relabel_configs: target.relabel_configs || [],
      metric_relabel_configs: target.metric_relabel_configs || []
    });
    setShowForm(true);
    if (target.relabel_configs?.length || target.metric_relabel_configs?.length) {
      setShowRelabelConfig(true);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('确定要删除这个目标吗？')) return;
    
    try {
      await TargetService.deleteTarget(id);
      onDataChange();
    } catch (error) {
      console.error('Error deleting target:', error);
      alert('删除失败，请重试');
    }
  };

  const addRelabelConfig = (type: 'relabel_configs' | 'metric_relabel_configs') => {
    const newConfig = {
      source_labels: [],
      target_label: '',
      regex: '(.*)',
      replacement: '${1}',
      action: 'replace' as const
    };
    
    setFormData({
      ...formData,
      [type]: [...formData[type], newConfig]
    });
  };

  const updateRelabelConfig = (
    type: 'relabel_configs' | 'metric_relabel_configs',
    index: number,
    field: string,
    value: any
  ) => {
    const configs = [...formData[type]];
    if (field === 'source_labels') {
      configs[index][field] = value.split(',').map((s: string) => s.trim()).filter((s: string) => s);
    } else {
      configs[index][field] = value;
    }
    setFormData({ ...formData, [type]: configs });
  };

  const removeRelabelConfig = (type: 'relabel_configs' | 'metric_relabel_configs', index: number) => {
    const configs = [...formData[type]];
    configs.splice(index, 1);
    setFormData({ ...formData, [type]: configs });
  };

  return (
    <div className="p-8">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold text-white mb-2">Target Management</h1>
          <p className="text-gray-400">Configure and manage Prometheus scrape targets</p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg flex items-center gap-2 transition-colors"
        >
          <Plus className="w-5 h-5" />
          Add Target
        </button>
      </div>

      {showForm && (
        <div className="bg-gray-800 rounded-xl p-6 border border-gray-700 mb-8">
          <h2 className="text-xl font-semibold text-white mb-6">
            {editingTarget ? 'Edit Target' : 'Add New Target'}
          </h2>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Job Name
                </label>
                <input
                  type="text"
                  value={formData.job_name}
                  onChange={(e) => setFormData({ ...formData, job_name: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="e.g., node-exporter"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Scrape Interval
                </label>
                <select
                  value={formData.scrape_interval}
                  onChange={(e) => setFormData({ ...formData, scrape_interval: e.target.value })}
                  className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="5s">5 seconds</option>
                  <option value="10s">10 seconds</option>
                  <option value="15s">15 seconds</option>
                  <option value="30s">30 seconds</option>
                  <option value="1m">1 minute</option>
                  <option value="5m">5 minutes</option>
                </select>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">
                Targets (comma-separated)
              </label>
              <input
                type="text"
                value={formData.targets}
                onChange={(e) => setFormData({ ...formData, targets: e.target.value })}
                className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="localhost:9100, localhost:9090"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">
                Metrics Path
              </label>
              <input
                type="text"
                value={formData.metrics_path}
                onChange={(e) => setFormData({ ...formData, metrics_path: e.target.value })}
                className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="/metrics"
                required
              />
            </div>

            {/* Relabel Configuration Toggle */}
            <div>
              <button
                type="button"
                onClick={() => setShowRelabelConfig(!showRelabelConfig)}
                className="flex items-center gap-2 text-gray-300 hover:text-white transition-colors"
              >
                {showRelabelConfig ? <ChevronDown className="w-4 h-4" /> : <ChevronRight className="w-4 h-4" />}
                <Settings className="w-4 h-4" />
                Relabel Configuration
              </button>
            </div>

            {showRelabelConfig && (
              <div className="space-y-6 border-t border-gray-700 pt-6">
                {/* Relabel Configs */}
                <div>
                  <div className="flex items-center justify-between mb-4">
                    <h3 className="text-lg font-medium text-white">Relabel Configs</h3>
                    <button
                      type="button"
                      onClick={() => addRelabelConfig('relabel_configs')}
                      className="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1 rounded text-sm transition-colors"
                    >
                      Add Rule
                    </button>
                  </div>
                  
                  {formData.relabel_configs.map((config, index) => (
                    <div key={index} className="bg-gray-700 rounded-lg p-4 mb-3">
                      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                        <div>
                          <label className="block text-xs text-gray-400 mb-1">Source Labels</label>
                          <input
                            type="text"
                            value={config.source_labels?.join(', ') || ''}
                            onChange={(e) => updateRelabelConfig('relabel_configs', index, 'source_labels', e.target.value)}
                            className="w-full bg-gray-600 border border-gray-500 rounded px-3 py-2 text-white text-sm"
                            placeholder="__meta_consul_service"
                          />
                        </div>
                        <div>
                          <label className="block text-xs text-gray-400 mb-1">Target Label</label>
                          <input
                            type="text"
                            value={config.target_label || ''}
                            onChange={(e) => updateRelabelConfig('relabel_configs', index, 'target_label', e.target.value)}
                            className="w-full bg-gray-600 border border-gray-500 rounded px-3 py-2 text-white text-sm"
                            placeholder="service"
                          />
                        </div>
                        <div>
                          <label className="block text-xs text-gray-400 mb-1">Action</label>
                          <select
                            value={config.action || 'replace'}
                            onChange={(e) => updateRelabelConfig('relabel_configs', index, 'action', e.target.value)}
                            className="w-full bg-gray-600 border border-gray-500 rounded px-3 py-2 text-white text-sm"
                          >
                            <option value="replace">replace</option>
                            <option value="keep">keep</option>
                            <option value="drop">drop</option>
                            <option value="hashmod">hashmod</option>
                            <option value="labelmap">labelmap</option>
                            <option value="labeldrop">labeldrop</option>
                            <option value="labelkeep">labelkeep</option>
                          </select>
                        </div>
                        <div>
                          <label className="block text-xs text-gray-400 mb-1">Regex</label>
                          <input
                            type="text"
                            value={config.regex || ''}
                            onChange={(e) => updateRelabelConfig('relabel_configs', index, 'regex', e.target.value)}
                            className="w-full bg-gray-600 border border-gray-500 rounded px-3 py-2 text-white text-sm"
                            placeholder="(.*)"
                          />
                        </div>
                        <div>
                          <label className="block text-xs text-gray-400 mb-1">Replacement</label>
                          <input
                            type="text"
                            value={config.replacement || ''}
                            onChange={(e) => updateRelabelConfig('relabel_configs', index, 'replacement', e.target.value)}
                            className="w-full bg-gray-600 border border-gray-500 rounded px-3 py-2 text-white text-sm"
                            placeholder="${1}"
                          />
                        </div>
                        <div className="flex items-end">
                          <button
                            type="button"
                            onClick={() => removeRelabelConfig('relabel_configs', index)}
                            className="bg-red-600 hover:bg-red-700 text-white px-3 py-2 rounded text-sm transition-colors"
                          >
                            Remove
                          </button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>

                {/* Metric Relabel Configs */}
                <div>
                  <div className="flex items-center justify-between mb-4">
                    <h3 className="text-lg font-medium text-white">Metric Relabel Configs</h3>
                    <button
                      type="button"
                      onClick={() => addRelabelConfig('metric_relabel_configs')}
                      className="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1 rounded text-sm transition-colors"
                    >
                      Add Rule
                    </button>
                  </div>
                  
                  {formData.metric_relabel_configs.map((config, index) => (
                    <div key={index} className="bg-gray-700 rounded-lg p-4 mb-3">
                      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                        <div>
                          <label className="block text-xs text-gray-400 mb-1">Source Labels</label>
                          <input
                            type="text"
                            value={config.source_labels?.join(', ') || ''}
                            onChange={(e) => updateRelabelConfig('metric_relabel_configs', index, 'source_labels', e.target.value)}
                            className="w-full bg-gray-600 border border-gray-500 rounded px-3 py-2 text-white text-sm"
                            placeholder="__name__"
                          />
                        </div>
                        <div>
                          <label className="block text-xs text-gray-400 mb-1">Target Label</label>
                          <input
                            type="text"
                            value={config.target_label || ''}
                            onChange={(e) => updateRelabelConfig('metric_relabel_configs', index, 'target_label', e.target.value)}
                            className="w-full bg-gray-600 border border-gray-500 rounded px-3 py-2 text-white text-sm"
                            placeholder="metric_name"
                          />
                        </div>
                        <div>
                          <label className="block text-xs text-gray-400 mb-1">Action</label>
                          <select
                            value={config.action || 'replace'}
                            onChange={(e) => updateRelabelConfig('metric_relabel_configs', index, 'action', e.target.value)}
                            className="w-full bg-gray-600 border border-gray-500 rounded px-3 py-2 text-white text-sm"
                          >
                            <option value="replace">replace</option>
                            <option value="keep">keep</option>
                            <option value="drop">drop</option>
                            <option value="hashmod">hashmod</option>
                            <option value="labelmap">labelmap</option>
                            <option value="labeldrop">labeldrop</option>
                            <option value="labelkeep">labelkeep</option>
                          </select>
                        </div>
                        <div>
                          <label className="block text-xs text-gray-400 mb-1">Regex</label>
                          <input
                            type="text"
                            value={config.regex || ''}
                            onChange={(e) => updateRelabelConfig('metric_relabel_configs', index, 'regex', e.target.value)}
                            className="w-full bg-gray-600 border border-gray-500 rounded px-3 py-2 text-white text-sm"
                            placeholder="(.*)"
                          />
                        </div>
                        <div>
                          <label className="block text-xs text-gray-400 mb-1">Replacement</label>
                          <input
                            type="text"
                            value={config.replacement || ''}
                            onChange={(e) => updateRelabelConfig('metric_relabel_configs', index, 'replacement', e.target.value)}
                            className="w-full bg-gray-600 border border-gray-500 rounded px-3 py-2 text-white text-sm"
                            placeholder="${1}"
                          />
                        </div>
                        <div className="flex items-end">
                          <button
                            type="button"
                            onClick={() => removeRelabelConfig('metric_relabel_configs', index)}
                            className="bg-red-600 hover:bg-red-700 text-white px-3 py-2 rounded text-sm transition-colors"
                          >
                            Remove
                          </button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            <div className="flex gap-4">
              <button
                type="submit"
                disabled={saving}
                className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg transition-colors"
              >
                {saving ? '保存中...' : (editingTarget ? 'Update Target' : 'Add Target')}
              </button>
              <button
                type="button"
                onClick={resetForm}
                className="bg-gray-600 hover:bg-gray-700 text-white px-6 py-3 rounded-lg transition-colors"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {targets.map((target) => (
          <div key={target.id} className="bg-gray-800 rounded-xl p-6 border border-gray-700 hover:border-gray-600 transition-colors">
            <div className="flex items-start justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-blue-500/20 rounded-lg">
                  <Database className="w-5 h-5 text-blue-400" />
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-white">{target.job_name}</h3>
                  <div className="flex items-center gap-2 mt-1">
                    <span className="px-2 py-1 bg-gray-700 text-gray-300 text-xs rounded">
                      {target.scrape_interval}
                    </span>
                    {target.relabel_configs && target.relabel_configs.length > 0 && (
                      <span className="px-2 py-1 bg-purple-600 text-white text-xs rounded">
                        {target.relabel_configs.length} Relabel Rules
                      </span>
                    )}
                    {target.metric_relabel_configs && target.metric_relabel_configs.length > 0 && (
                      <span className="px-2 py-1 bg-green-600 text-white text-xs rounded">
                        {target.metric_relabel_configs.length} Metric Rules
                      </span>
                    )}
                  </div>
                </div>
              </div>
              <div className="flex gap-2">
                <button
                  onClick={() => handleEdit(target)}
                  className="p-2 text-gray-400 hover:text-blue-400 hover:bg-blue-500/10 rounded-lg transition-colors"
                >
                  <Edit className="w-4 h-4" />
                </button>
                <button
                  onClick={() => handleDelete(target.id)}
                  className="p-2 text-gray-400 hover:text-red-400 hover:bg-red-500/10 rounded-lg transition-colors"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>

            <div className="space-y-2">
              <div>
                <p className="text-sm font-medium text-gray-300">Targets:</p>
                <p className="text-gray-400 text-sm">{target.targets.join(', ')}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-300">Metrics Path:</p>
                <p className="text-gray-400 text-sm">{target.metrics_path}</p>
              </div>
              {target.relabel_configs && target.relabel_configs.length > 0 && (
                <div>
                  <p className="text-sm font-medium text-gray-300">Relabel Rules:</p>
                  <div className="text-gray-400 text-sm">
                    {target.relabel_configs.map((rule, index) => (
                      <div key={index} className="bg-gray-700 rounded p-2 mt-1">
                        <span className="text-purple-400">{rule.action}</span>
                        {rule.source_labels && rule.source_labels.length > 0 && (
                          <span className="ml-2">from: {rule.source_labels.join(', ')}</span>
                        )}
                        {rule.target_label && (
                          <span className="ml-2">to: {rule.target_label}</span>
                        )}
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};