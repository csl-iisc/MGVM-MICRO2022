	.text
	.hsa_code_object_version 2,1
	.hsa_code_object_isa 8,0,3,"AMD","AMDGPU"
	.protected	compute_lj_force ; -- Begin function compute_lj_force
	.globl	compute_lj_force
	.p2align	8
	.type	compute_lj_force,@function
	.amdgpu_hsa_kernel compute_lj_force
compute_lj_force:                       ; @compute_lj_force
compute_lj_force$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 3
		granulated_wavefront_sgpr_count = 2
		priority = 0
		float_mode = 192
		priv = 0
		enable_dx10_clamp = 1
		debug_mode = 0
		enable_ieee_mode = 1
		enable_wgp_mode = 0
		enable_mem_ordered = 0
		enable_fwd_progress = 0
		enable_sgpr_private_segment_wave_byte_offset = 0
		user_sgpr_count = 8
		enable_trap_handler = 0
		enable_sgpr_workgroup_id_x = 1
		enable_sgpr_workgroup_id_y = 0
		enable_sgpr_workgroup_id_z = 0
		enable_sgpr_workgroup_info = 0
		enable_vgpr_workitem_id = 0
		enable_exception_msb = 0
		granulated_lds_size = 0
		enable_exception = 0
		enable_sgpr_private_segment_buffer = 1
		enable_sgpr_dispatch_ptr = 1
		enable_sgpr_queue_ptr = 0
		enable_sgpr_kernarg_segment_ptr = 1
		enable_sgpr_dispatch_id = 0
		enable_sgpr_flat_scratch_init = 0
		enable_sgpr_private_segment_size = 0
		enable_sgpr_grid_workgroup_count_x = 0
		enable_sgpr_grid_workgroup_count_y = 0
		enable_sgpr_grid_workgroup_count_z = 0
		enable_wavefront_size32 = 0
		enable_ordered_append_gds = 0
		private_element_size = 1
		is_ptr64 = 1
		is_dynamic_callstack = 0
		is_debug_enabled = 0
		is_xnack_enabled = 0
		workitem_private_segment_byte_size = 0
		workgroup_group_segment_byte_size = 0
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 104
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 22
		workitem_vgpr_count = 16
		reserved_vgpr_first = 0
		reserved_vgpr_count = 0
		reserved_sgpr_first = 0
		reserved_sgpr_count = 0
		debug_wavefront_private_segment_offset_sgpr = 0
		debug_private_segment_buffer_sgpr = 0
		kernarg_segment_alignment = 4
		group_segment_alignment = 4
		private_segment_alignment = 4
		wavefront_size = 6
		call_convention = -1
		runtime_loader_kernel_symbol = 0
	.end_amd_kernel_code_t
; %bb.0:
	s_load_dwordx4 s[0:3], s[6:7], 0x0
	s_load_dword s9, s[6:7], 0x10
	s_load_dwordx2 s[10:11], s[6:7], 0x30
	s_load_dword s4, s[4:5], 0x4
	s_mov_b32 s16, 0
	s_mov_b32 s17, s16
	s_mov_b32 s18, s16
	s_mov_b32 s19, s16
	s_waitcnt lgkmcnt(0)
	s_and_b32 s4, s4, 0xffff
	s_mul_i32 s8, s8, s4
	v_mov_b32_e32 v2, s16
	v_add_u32_e32 v0, vcc, s8, v0
	v_mov_b32_e32 v1, 0
	v_add_u32_e32 v0, vcc, s10, v0
	s_cmp_lt_i32 s9, 1
	v_mov_b32_e32 v3, s17
	v_mov_b32_e32 v4, s18
	v_mov_b32_e32 v5, s19
	s_cbranch_scc1 BB0_5
; %bb.1:
	v_lshlrev_b64 v[2:3], 4, v[0:1]
	v_mov_b32_e32 v4, s3
	v_add_u32_e32 v2, vcc, s2, v2
	v_addc_u32_e32 v3, vcc, v4, v3, vcc
	s_load_dwordx2 s[4:5], s[6:7], 0x18
	s_load_dwordx4 s[12:15], s[6:7], 0x20
	flat_load_dwordx4 v[6:9], v[2:3]
	v_mov_b32_e32 v2, s16
	v_mov_b32_e32 v3, s17
	v_mov_b32_e32 v4, s18
	s_waitcnt lgkmcnt(0)
	s_xor_b32 s6, s14, 0x80000000
	v_mov_b32_e32 v5, s19
	s_waitcnt vmcnt(0)
	v_mov_b32_e32 v9, v0
	s_branch BB0_3
BB0_2:                                  ;   in Loop: Header=BB0_3 Depth=1
	s_or_b64 exec, exec, s[10:11]
	s_add_i32 s9, s9, -1
	s_cmp_lg_u32 s9, 0
	v_add_u32_e32 v9, vcc, s15, v9
	s_cbranch_scc0 BB0_5
BB0_3:                                  ; =>This Inner Loop Header: Depth=1
	v_mov_b32_e32 v10, 0
	v_lshlrev_b64 v[10:11], 2, v[9:10]
	v_mov_b32_e32 v12, s5
	v_add_u32_e32 v10, vcc, s4, v10
	v_addc_u32_e32 v11, vcc, v12, v11, vcc
	flat_load_dword v10, v[10:11]
	v_mov_b32_e32 v12, s3
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_ashrrev_i32_e32 v11, 31, v10
	v_lshlrev_b64 v[10:11], 4, v[10:11]
	v_add_u32_e32 v10, vcc, s2, v10
	v_addc_u32_e32 v11, vcc, v12, v11, vcc
	flat_load_dwordx4 v[12:15], v[10:11]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_sub_f32_e32 v11, v7, v13
	v_sub_f32_e32 v12, v6, v12
	v_mul_f32_e32 v13, v11, v11
	v_sub_f32_e32 v10, v8, v14
	v_mac_f32_e32 v13, v12, v12
	v_mac_f32_e32 v13, v10, v10
	v_cmp_gt_f32_e32 vcc, s12, v13
	s_and_saveexec_b64 s[10:11], vcc
	s_cbranch_execz BB0_2
; %bb.4:                                ;   in Loop: Header=BB0_3 Depth=1
	v_rcp_f32_e32 v13, v13
	v_mov_b32_e32 v14, s6
	v_mul_f32_e32 v15, v13, v13
	v_mul_f32_e32 v15, v13, v15
	v_mul_f32_e32 v13, v13, v15
	v_mad_f32 v14, s13, v15, v14
	v_mul_f32_e32 v13, v13, v14
	v_mad_f32 v2, v12, v13, v2
	v_mad_f32 v3, v11, v13, v3
	v_mac_f32_e32 v4, v10, v13
	s_branch BB0_2
BB0_5:                                  ; %.loopexit
	v_lshlrev_b64 v[0:1], 4, v[0:1]
	v_mov_b32_e32 v6, s1
	v_add_u32_e32 v0, vcc, s0, v0
	v_addc_u32_e32 v1, vcc, v6, v1, vcc
	flat_store_dwordx4 v[0:1], v[2:5]
	s_endpgm
.Lfunc_end0:
	.size	compute_lj_force, .Lfunc_end0-compute_lj_force
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 400
; NumSgprs: 22
; NumVgprs: 16
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 0 bytes/workgroup (compile time only)
; SGPRBlocks: 2
; VGPRBlocks: 3
; NumSGPRsForWavesPerEU: 22
; NumVGPRsForWavesPerEU: 16
; Occupancy: 10
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.ident	"clang version 11.0.0 (/src/external/llvm-project/clang 0383ad1cfb0a8e05b0a020e8632400194628b243)"
	.section	".note.GNU-stack"
	.addrsig
	.amd_amdgpu_isa "amdgcn-amd-amdhsa--gfx803"
	.amd_amdgpu_hsa_metadata
---
Version:         [ 1, 0 ]
Kernels:
  - Name:            compute_lj_force
    SymbolName:      'compute_lj_force@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            force
        TypeName:        'float4*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            position
        TypeName:        'float4*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            neighCount
        TypeName:        int
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       I32
        AccQual:         Default
      - Name:            neighList
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            cutsq
        TypeName:        float
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       F32
        AccQual:         Default
      - Name:            lj1
        TypeName:        float
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       F32
        AccQual:         Default
      - Name:            lj2
        TypeName:        float
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       F32
        AccQual:         Default
      - Name:            inum
        TypeName:        int
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       I32
        AccQual:         Default
      - Size:            8
        Align:           8
        ValueKind:       HiddenGlobalOffsetX
        ValueType:       I64
      - Size:            8
        Align:           8
        ValueKind:       HiddenGlobalOffsetY
        ValueType:       I64
      - Size:            8
        Align:           8
        ValueKind:       HiddenGlobalOffsetZ
        ValueType:       I64
      - Size:            8
        Align:           8
        ValueKind:       HiddenNone
        ValueType:       I8
        AddrSpaceQual:   Global
      - Size:            8
        Align:           8
        ValueKind:       HiddenNone
        ValueType:       I8
        AddrSpaceQual:   Global
      - Size:            8
        Align:           8
        ValueKind:       HiddenNone
        ValueType:       I8
        AddrSpaceQual:   Global
      - Size:            8
        Align:           8
        ValueKind:       HiddenMultiGridSyncArg
        ValueType:       I8
        AddrSpaceQual:   Global
    CodeProps:
      KernargSegmentSize: 104
      GroupSegmentFixedSize: 0
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        22
      NumVGPRs:        16
      MaxFlatWorkGroupSize: 256
...

	.end_amd_amdgpu_hsa_metadata
