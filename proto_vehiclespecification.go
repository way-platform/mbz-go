package mbz

import (
	"github.com/way-platform/mbz-go/api/vehiclespecificationfleetv1"
	mbzv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mbz/v1"
)

func vehicleDataToProto(
	vehicleData *vehiclespecificationfleetv1.VehicleData,
) *mbzv1.VehicleSpecification {
	output := &mbzv1.VehicleSpecification{}
	if vehicleData == nil {
		return output
	}
	// Basic vehicle identification fields
	if vehicleData.Model != "" {
		output.SetModel(vehicleData.Model)
	}
	if vehicleData.ModelName != "" {
		output.SetModelName(vehicleData.ModelName)
	}
	if vehicleData.ModelYear != "" {
		output.SetModelYear(vehicleData.ModelYear)
	}
	if vehicleData.ModelYearYear != "" {
		output.SetModelYearYear(vehicleData.ModelYearYear)
	}
	if vehicleData.LongType != "" {
		output.SetLongType(vehicleData.LongType)
	}
	if vehicleData.LongTypeTechnical != "" {
		output.SetLongTypeTechnical(vehicleData.LongTypeTechnical)
	}
	if vehicleData.ShortType != "" {
		output.SetShortType(vehicleData.ShortType)
	}
	if vehicleData.EngineCode != "" {
		output.SetEngineCode(vehicleData.EngineCode)
	}
	if vehicleData.Fin11 != "" {
		output.SetFin11(vehicleData.Fin11)
	}
	if vehicleData.Vin11 != "" {
		output.SetVin11(vehicleData.Vin11)
	}
	if vehicleData.Baumuster != "" {
		output.SetBaumuster(vehicleData.Baumuster)
	}
	if vehicleData.Countrynumber != "" {
		output.SetCountrynumber(vehicleData.Countrynumber)
	}
	if vehicleData.Function != "" {
		output.SetFunction(vehicleData.Function)
	}
	if vehicleData.Status != "" {
		output.SetStatus(vehicleData.Status)
	}
	if vehicleData.Finaldate != "" {
		output.SetFinaldate(vehicleData.Finaldate)
	}
	if vehicleData.Shippingdate != "" {
		output.SetShippingdate(vehicleData.Shippingdate)
	}
	if vehicleData.Nst != "" {
		output.SetNst(vehicleData.Nst)
	}
	if vehicleData.Typeclass != "" {
		output.SetTypeclass(vehicleData.Typeclass)
	}
	if vehicleData.Allwheel != "" {
		output.SetAllwheel(vehicleData.Allwheel)
	}
	if vehicleData.Metallic != "" {
		output.SetMetallic(vehicleData.Metallic)
	}
	if vehicleData.Leather != "" {
		output.SetLeather(vehicleData.Leather)
	}
	if vehicleData.VehicleType != "" {
		output.SetVehicleType(string(vehicleData.VehicleType))
	}
	if vehicleData.Poolflag != "" {
		output.SetPoolflag(string(vehicleData.Poolflag))
	}

	// CodeText fields
	if ct := convertCodeText(vehicleData.Brand); ct != nil {
		output.SetBrand(ct)
	}
	if ct := convertCodeText(vehicleData.Fuel); ct != nil {
		output.SetFuel(ct)
	}
	if ct := convertCodeText(vehicleData.Emissionstandard); ct != nil {
		output.SetEmissionstandard(ct)
	}
	if ct := convertCodeText(vehicleData.Body); ct != nil {
		output.SetBody(ct)
	}
	if ct := convertCodeText(vehicleData.Cabin); ct != nil {
		output.SetCabin(ct)
	}
	if ct := convertCodeText(vehicleData.Enginetype); ct != nil {
		output.SetEnginetype(ct)
	}
	if ct := convertCodeText(vehicleData.Line); ct != nil {
		output.SetLine(ct)
	}
	if ct := convertCodeText(vehicleData.Modelrange); ct != nil {
		output.SetModelrange(ct)
	}
	if ct := convertCodeText(vehicleData.Productgroup); ct != nil {
		output.SetProductgroup(ct)
	}
	if ct := convertCodeText(vehicleData.Subcategory); ct != nil {
		output.SetSubcategory(ct)
	}
	if ct := convertCodeText(vehicleData.Suspensiontype); ct != nil {
		output.SetSuspensiontype(ct)
	}
	if ct := convertCodeText(vehicleData.Transmissiontype); ct != nil {
		output.SetTransmissiontype(ct)
	}

	// Numeric fields
	if vehicleData.Branch != nil && *vehicleData.Branch > 0 {
		output.SetBranch(*vehicleData.Branch)
	}
	if vehicleData.Cylindercapacity != nil && *vehicleData.Cylindercapacity > 0 {
		output.SetCylinderCapacityCc(*vehicleData.Cylindercapacity)
	}
	if vehicleData.Powerkw != nil && *vehicleData.Powerkw > 0 {
		output.SetPowerKw(*vehicleData.Powerkw)
	}
	if vehicleData.Powerps != nil && *vehicleData.Powerps > 0 {
		output.SetPowerPs(*vehicleData.Powerps)
	}
	if vehicleData.Numberofdoors != nil && *vehicleData.Numberofdoors > 0 {
		output.SetDoorCount(*vehicleData.Numberofdoors)
	}
	if vehicleData.Numberofseats != nil && *vehicleData.Numberofseats > 0 {
		output.SetSeatCount(*vehicleData.Numberofseats)
	}
	if vehicleData.Wheelbase != nil && *vehicleData.Wheelbase > 0 {
		output.SetWheelbaseMm(*vehicleData.Wheelbase)
	}
	if vehicleData.Wheelform != "" {
		output.SetWheelForm(vehicleData.Wheelform)
	}

	// Complex nested structures
	if w := convertWeight(vehicleData.Weight); w != nil {
		output.SetWeight(w)
	}
	if e := convertEngine(vehicleData.PrimaryEngine); e != nil {
		output.SetPrimaryEngine(e)
	}
	if e := convertEngine(vehicleData.SecondaryEngine); e != nil {
		output.SetSecondaryEngine(e)
	}
	if ci := convertCabinInfo(vehicleData.CabinInfo); ci != nil {
		output.SetCabinInfo(ci)
	}
	if e := convertEmico(vehicleData.Emico); e != nil {
		output.SetEmico(e)
	}
	if e := convertEmission(vehicleData.Emission); e != nil {
		output.SetEmission(e)
	}
	if fc := convertFuelConsumption(vehicleData.Fuelconsumption); fc != nil {
		output.SetFuelConsumption(fc)
	}
	if np := convertNewPrice(vehicleData.Newprice); np != nil {
		output.SetNewPrice(np)
	}
	if sd := convertSalesDescription(vehicleData.SalesDescription); sd != nil {
		output.SetSalesDescription(sd)
	}
	if pbo := convertPackageBasedOptions(vehicleData.PackageBasedOptions); pbo != nil {
		output.SetPackageBasedOptions(pbo)
	}
	if tk := convertTypKey(vehicleData.Typekey); tk != nil {
		output.SetTypeKey(tk)
	}
	if le := convertLegacyEquipment(vehicleData.Paint); le != nil {
		output.SetPaint(le)
	}
	if le := convertLegacyEquipment(vehicleData.Paint2); le != nil {
		output.SetPaint2(le)
	}
	if le := convertLegacyEquipment(vehicleData.Upholstery); le != nil {
		output.SetUpholstery(le)
	}

	// Arrays
	if len(vehicleData.Options) > 0 {
		options := make([]*mbzv1.VehicleSpecification_LegacyEquipment, 0, len(vehicleData.Options))
		for _, opt := range vehicleData.Options {
			if o := convertLegacyEquipment(&opt); o != nil {
				options = append(options, o)
			}
		}
		if len(options) > 0 {
			output.SetOptions(options)
		}
	}
	if len(vehicleData.TechnicalData) > 0 {
		technicalData := make([]*mbzv1.VehicleSpecification_TechnicalData, 0, len(vehicleData.TechnicalData))
		for _, td := range vehicleData.TechnicalData {
			if t := convertTechnicalData(&td); t != nil {
				technicalData = append(technicalData, t)
			}
		}
		if len(technicalData) > 0 {
			output.SetTechnicalData(technicalData)
		}
	}

	return output
}

func convertCodeText(ct *vehiclespecificationfleetv1.CodeText) *mbzv1.VehicleSpecification_CodeText {
	if ct == nil {
		return nil
	}
	// Only create CodeText if at least one field is non-empty
	if ct.Code == "" && ct.Text == "" {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_CodeText{}
	if ct.Code != "" {
		proto.SetCode(ct.Code)
	}
	if ct.Text != "" {
		proto.SetText(ct.Text)
	}
	return proto
}

func convertWeight(w *vehiclespecificationfleetv1.Weight) *mbzv1.VehicleSpecification_Weight {
	if w == nil {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_Weight{}
	hasFields := false
	if w.VehicleMassKg != nil && *w.VehicleMassKg > 0 {
		proto.SetVehicleMassKg(*w.VehicleMassKg)
		hasFields = true
	}
	if w.Total != nil && *w.Total > 0 {
		proto.SetTotalKg(*w.Total)
		hasFields = true
	}
	if w.Payload != nil && *w.Payload > 0 {
		proto.SetPayloadKg(int32(*w.Payload))
		hasFields = true
	}
	if w.Extended != nil && *w.Extended > 0 {
		proto.SetExtendedKg(*w.Extended)
		hasFields = true
	}
	if w.Axis1 != nil && *w.Axis1 > 0 {
		proto.SetAxis1Kg(*w.Axis1)
		hasFields = true
	}
	if w.Axis2 != nil && *w.Axis2 > 0 {
		proto.SetAxis2Kg(*w.Axis2)
		hasFields = true
	}
	if w.Axis3 != nil && *w.Axis3 > 0 {
		proto.SetAxis3Kg(*w.Axis3)
		hasFields = true
	}
	if w.Axis4 != nil && *w.Axis4 > 0 {
		proto.SetAxis4Kg(*w.Axis4)
		hasFields = true
	}
	if w.Axis5 != nil && *w.Axis5 > 0 {
		proto.SetAxis5Kg(*w.Axis5)
		hasFields = true
	}
	if w.Axis6 != nil && *w.Axis6 > 0 {
		proto.SetAxis6Kg(*w.Axis6)
		hasFields = true
	}
	if w.Axislast != nil && *w.Axislast > 0 {
		proto.SetAxisLastKg(*w.Axislast)
		hasFields = true
	}
	if !hasFields {
		return nil
	}
	return proto
}

func convertEngine(e *vehiclespecificationfleetv1.Engine) *mbzv1.VehicleSpecification_Engine {
	if e == nil {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_Engine{}
	hasFields := false
	if b := convertEngineBattery(e.Battery); b != nil {
		proto.SetBattery(b)
		hasFields = true
	}
	if e.CylinderCapacity != nil && *e.CylinderCapacity > 0 {
		proto.SetCylinderCapacityCc(*e.CylinderCapacity)
		hasFields = true
	}
	if em := convertEngineEmission(e.Emission); em != nil {
		proto.SetEmission(em)
		hasFields = true
	}
	if fc := convertEngineFuelConsumption(e.FuelConsumption); fc != nil {
		proto.SetFuelConsumption(fc)
		hasFields = true
	}
	if e.FuelTankCapacity != nil && *e.FuelTankCapacity > 0 {
		proto.SetFuelTankCapacityL(*e.FuelTankCapacity)
		hasFields = true
	}
	if ft := convertCodeText(e.FuelType); ft != nil {
		proto.SetFuelType(ft)
		hasFields = true
	}
	if p := convertEnginePower(e.Power); p != nil {
		proto.SetPower(p)
		hasFields = true
	}
	if e.Range != nil && *e.Range > 0 {
		proto.SetRangeKm(*e.Range)
		hasFields = true
	}
	if e.Torque != nil && *e.Torque > 0 {
		proto.SetTorqueNm(*e.Torque)
		hasFields = true
	}
	if !hasFields {
		return nil
	}
	return proto
}

func convertEngineBattery(b *vehiclespecificationfleetv1.EngineBattery) *mbzv1.VehicleSpecification_Engine_Battery {
	if b == nil {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_Engine_Battery{}
	hasFields := false
	if b.Capacity != nil && *b.Capacity > 0 {
		proto.SetCapacityKwh(*b.Capacity)
		hasFields = true
	}
	if b.ChargingTime != nil && *b.ChargingTime > 0 {
		proto.SetChargingTimeMin(*b.ChargingTime)
		hasFields = true
	}
	if b.Type != "" {
		proto.SetType(b.Type)
		hasFields = true
	}
	if !hasFields {
		return nil
	}
	return proto
}

func convertEngineEmission(e *vehiclespecificationfleetv1.EngineEmission) *mbzv1.VehicleSpecification_Engine_EngineEmission {
	if e == nil {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_Engine_EngineEmission{}
	hasFields := false
	if e.City != nil {
		proto.SetCity(*e.City)
		hasFields = true
	}
	if e.Combined != nil {
		proto.SetCombined(*e.Combined)
		hasFields = true
	}
	if e.Overland != nil {
		proto.SetOverland(*e.Overland)
		hasFields = true
	}
	if e.Unit != "" {
		proto.SetUnit(e.Unit)
		hasFields = true
	}
	if !hasFields {
		return nil
	}
	return proto
}

func convertEngineFuelConsumption(fc *vehiclespecificationfleetv1.EngineFuelConsumption) *mbzv1.VehicleSpecification_Engine_EngineFuelConsumption {
	if fc == nil {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_Engine_EngineFuelConsumption{}
	hasFields := false
	if fc.City != nil {
		proto.SetCity(*fc.City)
		hasFields = true
	}
	if fc.Combined != nil {
		proto.SetCombined(*fc.Combined)
		hasFields = true
	}
	if fc.Overland != nil {
		proto.SetOverland(*fc.Overland)
		hasFields = true
	}
	if fc.Unit != "" {
		proto.SetUnit(fc.Unit)
		hasFields = true
	}
	if !hasFields {
		return nil
	}
	return proto
}

func convertEnginePower(p *vehiclespecificationfleetv1.EnginePower) *mbzv1.VehicleSpecification_Engine_Power {
	if p == nil {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_Engine_Power{}
	hasFields := false
	if p.Kw != nil && *p.Kw > 0 {
		proto.SetKw(*p.Kw)
		hasFields = true
	}
	if p.Ps != nil && *p.Ps > 0 {
		proto.SetPs(*p.Ps)
		hasFields = true
	}
	if !hasFields {
		return nil
	}
	return proto
}

func convertCabinInfo(ci *vehiclespecificationfleetv1.CabinInfo) *mbzv1.VehicleSpecification_CabinInfo {
	if ci == nil {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_CabinInfo{}
	hasFields := false
	if ct := convertCodeText(ci.CabinBase); ct != nil {
		proto.SetCabinBase(ct)
		hasFields = true
	}
	if ct := convertCodeText(ci.CabinBottom); ct != nil {
		proto.SetCabinBottom(ct)
		hasFields = true
	}
	if ct := convertCodeText(ci.CabinCombined); ct != nil {
		proto.SetCabinCombined(ct)
		hasFields = true
	}
	if ct := convertCodeText(ci.CabinType); ct != nil {
		proto.SetCabinType(ct)
		hasFields = true
	}
	if ct := convertCodeText(ci.CabinWidth); ct != nil {
		proto.SetCabinWidth(ct)
		hasFields = true
	}
	if !hasFields {
		return nil
	}
	return proto
}

func convertEmico(e *vehiclespecificationfleetv1.Emico) *mbzv1.VehicleSpecification_Emico {
	if e == nil {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_Emico{}
	hasFields := false
	if e.ApprovedUsage != "" {
		proto.SetApprovedUsage(e.ApprovedUsage)
		hasFields = true
	}
	if e.DriveConcept != "" {
		proto.SetDriveConcept(e.DriveConcept)
		hasFields = true
	}
	if e.LegalContext != "" {
		proto.SetLegalContext(e.LegalContext)
		hasFields = true
	}
	if e.RawData != "" {
		proto.SetRawData(e.RawData)
		hasFields = true
	}
	if e.RuleSetVersionID != "" {
		proto.SetRuleSetVersionId(e.RuleSetVersionID)
		hasFields = true
	}
	if e.TestProcedure != "" {
		proto.SetTestProcedure(e.TestProcedure)
		hasFields = true
	}
	if !hasFields {
		return nil
	}
	return proto
}

func convertEmission(e *vehiclespecificationfleetv1.Emission) *mbzv1.VehicleSpecification_Emission {
	if e == nil {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_Emission{}
	hasFields := false
	if e.Combined != nil {
		proto.SetCombined(*e.Combined)
		hasFields = true
	}
	if e.Directive != "" {
		proto.SetDirective(e.Directive)
		hasFields = true
	}
	if e.In != nil {
		proto.SetIn(*e.In)
		hasFields = true
	}
	if e.Out != nil {
		proto.SetOut(*e.Out)
		hasFields = true
	}
	if !hasFields {
		return nil
	}
	return proto
}

func convertFuelConsumption(fc *vehiclespecificationfleetv1.FuelConsumption) *mbzv1.VehicleSpecification_FuelConsumption {
	if fc == nil {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_FuelConsumption{}
	hasFields := false
	if fc.Combined != nil {
		proto.SetCombined(*fc.Combined)
		hasFields = true
	}
	if fc.In != nil {
		proto.SetIn(*fc.In)
		hasFields = true
	}
	if fc.Out != nil {
		proto.SetOut(*fc.Out)
		hasFields = true
	}
	if !hasFields {
		return nil
	}
	return proto
}

func convertNewPrice(np *vehiclespecificationfleetv1.NewPrice) *mbzv1.VehicleSpecification_NewPrice {
	if np == nil {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_NewPrice{}
	hasFields := false
	if np.Grossprice != nil {
		proto.SetGrossPrice(*np.Grossprice)
		hasFields = true
	}
	if np.Netprice != nil {
		proto.SetNetPrice(*np.Netprice)
		hasFields = true
	}
	if np.Tax != nil {
		proto.SetTax(*np.Tax)
		hasFields = true
	}
	if !hasFields {
		return nil
	}
	return proto
}

func convertSalesDescription(sd *vehiclespecificationfleetv1.SalesDescription) *mbzv1.VehicleSpecification_SalesDescription {
	if sd == nil {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_SalesDescription{}
	hasFields := false
	if ct := convertCodeText(sd.Allwheel); ct != nil {
		proto.SetAllwheel(ct)
		hasFields = true
	}
	if ct := convertCodeText(sd.Body); ct != nil {
		proto.SetBody(ct)
		hasFields = true
	}
	if ct := convertCodeText(sd.Extension); ct != nil {
		proto.SetExtension(ct)
		hasFields = true
	}
	if ct := convertCodeText(sd.LengthType); ct != nil {
		proto.SetLengthType(ct)
		hasFields = true
	}
	if ct := convertCodeText(sd.Line); ct != nil {
		proto.SetLine(ct)
		hasFields = true
	}
	if ct := convertCodeText(sd.Model); ct != nil {
		proto.SetModel(ct)
		hasFields = true
	}
	if ct := convertCodeText(sd.Style); ct != nil {
		proto.SetStyle(ct)
		hasFields = true
	}
	if ct := convertCodeText(sd.Type); ct != nil {
		proto.SetType(ct)
		hasFields = true
	}
	if !hasFields {
		return nil
	}
	return proto
}

func convertPackageBasedOptions(pbo *vehiclespecificationfleetv1.PackageBasedOptions) *mbzv1.VehicleSpecification_PackageBasedOptions {
	if pbo == nil {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_PackageBasedOptions{}
	hasFields := false
	if len(pbo.OptionPackages) > 0 {
		packages := make([]*mbzv1.VehicleSpecification_OptionPackage, 0, len(pbo.OptionPackages))
		for _, pkg := range pbo.OptionPackages {
			if p := convertOptionPackage(&pkg); p != nil {
				packages = append(packages, p)
			}
		}
		if len(packages) > 0 {
			proto.SetOptionPackages(packages)
			hasFields = true
		}
	}
	if len(pbo.Options) > 0 {
		options := make([]*mbzv1.VehicleSpecification_LegacyEquipment, 0, len(pbo.Options))
		for _, opt := range pbo.Options {
			if o := convertLegacyEquipment(&opt); o != nil {
				options = append(options, o)
			}
		}
		if len(options) > 0 {
			proto.SetOptions(options)
			hasFields = true
		}
	}
	if !hasFields {
		return nil
	}
	return proto
}

func convertOptionPackage(op *vehiclespecificationfleetv1.OptionPackage) *mbzv1.VehicleSpecification_OptionPackage {
	if op == nil {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_OptionPackage{}
	hasFields := false
	if op.Code != "" {
		proto.SetCode(op.Code)
		hasFields = true
	}
	if op.Text != "" {
		proto.SetText(op.Text)
		hasFields = true
	}
	if len(op.Options) > 0 {
		options := make([]*mbzv1.VehicleSpecification_LegacyEquipment, 0, len(op.Options))
		for _, opt := range op.Options {
			if o := convertLegacyEquipment(&opt); o != nil {
				options = append(options, o)
			}
		}
		if len(options) > 0 {
			proto.SetOptions(options)
			hasFields = true
		}
	}
	if !hasFields {
		return nil
	}
	return proto
}

func convertTechnicalData(td *vehiclespecificationfleetv1.TechnicalData) *mbzv1.VehicleSpecification_TechnicalData {
	if td == nil {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_TechnicalData{}
	hasFields := false
	if td.ID != "" {
		proto.SetId(td.ID)
		hasFields = true
	}
	if td.Text != "" {
		proto.SetText(td.Text)
		hasFields = true
	}
	if td.Unit != "" {
		proto.SetUnit(td.Unit)
		hasFields = true
	}
	if td.Value != "" {
		proto.SetValue(td.Value)
		hasFields = true
	}
	if !hasFields {
		return nil
	}
	return proto
}

func convertTypKey(tk *vehiclespecificationfleetv1.TypKey) *mbzv1.VehicleSpecification_TypKey {
	if tk == nil {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_TypKey{}
	hasFields := false
	if tk.Hsn != "" {
		proto.SetHsn(tk.Hsn)
		hasFields = true
	}
	if tk.Tsn != "" {
		proto.SetTsn(tk.Tsn)
		hasFields = true
	}
	if tk.Vvs != "" {
		proto.SetVvs(tk.Vvs)
		hasFields = true
	}
	if !hasFields {
		return nil
	}
	return proto
}

func convertLegacyEquipment(le *vehiclespecificationfleetv1.LegacyEquipment) *mbzv1.VehicleSpecification_LegacyEquipment {
	if le == nil {
		return nil
	}
	proto := &mbzv1.VehicleSpecification_LegacyEquipment{}
	hasFields := false
	if le.Code != "" {
		proto.SetCode(le.Code)
		hasFields = true
	}
	if le.CodeType != "" {
		proto.SetCodeType(le.CodeType)
		hasFields = true
	}
	if le.Description != "" {
		proto.SetDescription(le.Description)
		hasFields = true
	}
	if le.Group != "" {
		proto.SetGroup(le.Group)
		hasFields = true
	}
	if le.Type != "" {
		proto.SetType(le.Type)
		hasFields = true
	}
	if !hasFields {
		return nil
	}
	return proto
}
