import React, { useEffect, useState } from "react";
import axios from "axios";

type Condition = {
  ruleType: string;
  ruleMethod: string;
  ruleDefinition: string;
};

type RuleInput = {
  ruleID: string;
  action: string;
  category: string;
  conditions: Condition[];
  applications: string[];
};

type AppOption = {
  application_id: string;
  application_name: string;
};

interface EditRuleModalProps {
  ruleID: string;
  onClose: () => void;
  onSuccess: () => void;
}

const backendUrl = import.meta.env.VITE_BACKEND_URL;
const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDY4NTMyMTMsInJvbGUiOiJzdXBlcl9hZG1pbiIsInVzZXJfaWQiOiJiNGM1ZjI0OC1iOTE3LTQyNDMtYjE0ZS1kNmI4NWQ2NzZjODgifQ.PfJJkKGbHXHQ9xTirmBoE-VHM3Zp8xkZKIfdYWI8QWI";

const validRuleTypes = ["REQUEST_HEADERS", "ARGS", "REQUEST_METHOD", "REMOTE_ADDR"];
const validRuleMethods = ["regex", "contains", "ipMatch"];
const validActions = ["deny", "log", "drop", "pass", "redirect", "status:403"];

const EditRuleModal: React.FC<EditRuleModalProps> = ({ ruleID, onClose, onSuccess }) => {
  const [ruleInput, setRuleInput] = useState<RuleInput | null>(null);
  const [availableApps, setAvailableApps] = useState<AppOption[]>([]);
  const [preview, setPreview] = useState<string>("");

  useEffect(() => {
    const fetchRule = async () => {
      try {
        const res = await axios.get(`${backendUrl}/rule/${ruleID}`, {
          headers: { Authorization: token },
        });
        const rule = res.data.rule;
        const rule_def= res.data.rule_definition
        const applications = res.data.applications.map((app: any) => ({
          application_id: app.application_id,
          application_name: app.application_name,
        }));
        setRuleInput({
          ruleID: rule.rule_id,
          action: rule.action,
          category: rule.category,
          conditions: rule_def.map((c: any) => ({
            ruleType: c.rule_type,
            ruleMethod: c.rule_method,
            ruleDefinition: c.rule_definition,
          })),
          applications: applications.map((app: any) => app.application_id),
        });
        console.log(ruleInput)
      } catch (err) {
        console.error("Failed to fetch rule", err);
        alert("Could not load rule.");
        onClose();
      }
    };

    const fetchApps = async () => {
      try {
        const res = await axios.get(`${backendUrl}/application`, {
          headers: { Authorization: token },
        });
        const apps = res.data.applications.map((app: any) => ({
          application_id: app.application_id,
          application_name: app.application_name,
        }));
        setAvailableApps(apps);
      } catch (err) {
        console.error("Failed to fetch applications", err);
      }
    };

    fetchRule();
    fetchApps();
  }, [ruleID]);

  useEffect(() => {
    if (ruleInput) generateRule(ruleInput);
  }, [ruleInput]);

  const generateRule = (input: RuleInput) => {
    const { ruleID, action, category, conditions } = input;
    if (conditions.length === 0) return;

    let ruleText = "";
    conditions.forEach((cond, i) => {
      const prefix = i === 0 ? "SecRule" : "    SecRule";
      const chain = i < conditions.length - 1 ? `"chain"` : "";
      const firstLine = i === 0
        ? `"id:${ruleID},phase:2,${action},msg:'${category}'${conditions.length > 1 ? ",chain" : ""}"`
        : chain;
      ruleText += `${prefix} ${cond.ruleType} "@${cond.ruleMethod} ${cond.ruleDefinition}" ${firstLine}\n`;
    });

    setPreview(ruleText.trim());
  };

  const updateCondition = (index: number, field: keyof Condition, value: string) => {
    if (!ruleInput) return;
    const updated = [...ruleInput.conditions];
    updated[index][field] = value;
    setRuleInput({ ...ruleInput, conditions: updated });
  };

  const addCondition = () => {
    if (!ruleInput) return;
    setRuleInput({
      ...ruleInput,
      conditions: [...ruleInput.conditions, { ruleType: "", ruleMethod: "", ruleDefinition: "" }],
    });
  };

  const handleAppAdd = (appId: string) => {
    if (!ruleInput || ruleInput.applications.includes(appId)) return;
    setRuleInput({ ...ruleInput, applications: [...ruleInput.applications, appId] });
  };

  const handleAppRemove = (appId: string) => {
    if (!ruleInput) return;
    setRuleInput({
      ...ruleInput,
      applications: ruleInput.applications.filter(a => a !== appId),
    });
  };

  const saveRule = async () => {
    if (!ruleInput) return;
    console.log(ruleInput.applications)
    try {
      await axios.put(`${backendUrl}/rule/update/${ruleInput.ruleID}`, {
        action: ruleInput.action,
        category: ruleInput.category,
        is_active: true,
        conditions: ruleInput.conditions.map(cond => ({
          rule_type: cond.ruleType,
          rule_method: cond.ruleMethod,
          rule_definition: cond.ruleDefinition,
        })),
        application_ids: ruleInput.applications,
      }, {
        headers: { Authorization: token },
      });

      alert("Rule updated successfully.");
      onSuccess();
      onClose();
    } catch (err) {
      console.error("Failed to update rule", err);
      alert("Update failed.");
    }
  };

  if (!ruleInput) {
    return <div className="p-4">Loading rule data...</div>;
  }

  const unselectedApps = availableApps.filter(app => !ruleInput.applications.includes(app.application_id));

  return (
    <div className="p-4 max-w-4xl mx-auto space-y-4">
      <h2 className="text-xl font-bold">Edit Custom Rule</h2>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <select
          className="border p-2"
          value={ruleInput.action}
          onChange={(e) => setRuleInput({ ...ruleInput, action: e.target.value })}
        >
          <option value="">Select Action</option>
          {validActions.map(a => <option key={a} value={a}>{a}</option>)}
        </select>

        <input
          className="border p-2"
          placeholder="Category"
          value={ruleInput.category}
          onChange={(e) => setRuleInput({ ...ruleInput, category: e.target.value })}
        />
      </div>

      {/* Conditions */}
      <div>
        <h3 className="font-semibold">Conditions</h3>
        {ruleInput.conditions.map((cond, i) => (
          <div key={i} className="grid grid-cols-1 md:grid-cols-3 gap-2 mb-2">
            <select
              value={cond.ruleType}
              className="border p-2"
              onChange={(e) => updateCondition(i, "ruleType", e.target.value)}
            >
              <option value="">Type</option>
              {validRuleTypes.map(t => <option key={t} value={t}>{t}</option>)}
            </select>

            <select
              value={cond.ruleMethod}
              className="border p-2"
              onChange={(e) => updateCondition(i, "ruleMethod", e.target.value)}
            >
              <option value="">Method</option>
              {validRuleMethods.map(m => <option key={m} value={m}>{m}</option>)}
            </select>

            <input
              value={cond.ruleDefinition}
              className="border p-2"
              onChange={(e) => updateCondition(i, "ruleDefinition", e.target.value)}
              placeholder="Definition"
            />
          </div>
        ))}
        <button className="bg-blue-500 text-white px-4 py-1 rounded" onClick={addCondition}>+ Add Condition</button>
      </div>

      {/* Applications */}
      <div>
        <label className="block font-semibold mb-2">Target Applications</label>
        <select
          className="border p-2 w-full"
          onChange={(e) => {
            const value = e.target.value;
            if (value) handleAppAdd(value);
            e.target.value = "";
          }}
        >
          <option value="">-- Select Application --</option>
          {unselectedApps.map(app => (
            <option key={app.application_id} value={app.application_id}>
              {app.application_name}
            </option>
          ))}
        </select>

        <div className="flex flex-wrap gap-2 mt-3">
          {ruleInput.applications.map(appId => {
            const app = availableApps.find(a => a.application_id === appId);
            return (
              <span key={appId}
                className="bg-green-200 text-green-900 px-3 py-1 rounded-full cursor-pointer hover:bg-red-200"
                onClick={() => handleAppRemove(appId)}
                title="Click to remove"
              >
                {app?.application_name || appId}
              </span>
            );
          })}
        </div>
      </div>

      {/* Preview & Save
      <div>
        <label className="font-semibold block mb-1">Rule Preview</label>
        <textarea readOnly className="w-full border p-2 h-40 bg-gray-100">{preview}</textarea>
      </div> */}

      <div className="flex gap-4">
        <button className="bg-green-600 text-white px-4 py-2 rounded" onClick={saveRule}>Save Changes</button>
        <button className="bg-gray-400 text-white px-4 py-2 rounded" onClick={onClose}>Cancel</button>
      </div>
    </div>
  );
};

export default EditRuleModal;
